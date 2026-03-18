package runner

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	log "github.com/s00500/env_logger"
)

type CommandConfig struct {
	Dir        string        // working dir (default: cwd)
	WatchExts  []string      // file extensions to watch (default: [".templ"])
	IgnoreDirs []string      // dirs to ignore (default: .git, vendor, node_modules, tmp)
	Command    HookConfig    // the main command to run on file change
	PostHooks  []HookConfig  // commands to run after successful main command
	Master     bool          // touch sync file after each run
	Sync       bool          // watch sync file for reload signals
	Interactive bool         // manual reload via Enter key (default: true)
	Debounce   time.Duration // default: 200ms
}

type CommandRunner struct {
	config CommandConfig
}

func NewCommandRunner(config CommandConfig) *CommandRunner {
	if config.Dir == "" {
		config.Dir, _ = os.Getwd()
	}
	if len(config.WatchExts) == 0 {
		config.WatchExts = []string{".templ"}
	}
	if len(config.IgnoreDirs) == 0 {
		config.IgnoreDirs = []string{".git", "vendor", "node_modules", "tmp"}
	}
	if config.Debounce == 0 {
		config.Debounce = 200 * time.Millisecond
	}
	return &CommandRunner{config: config}
}

func (cr *CommandRunner) Run(ctx context.Context) error {
	ctx, cancel := signal.NotifyContext(ctx, notifySignals()...)
	defer cancel()

	reload := make(chan struct{}, 1)

	w, err := newWatcher(cr.config.Dir, cr.config.WatchExts, cr.config.IgnoreDirs, cr.config.Debounce)
	if err != nil {
		return fmt.Errorf("watcher: %w", err)
	}
	defer w.close()
	go w.watch(ctx, reload)

	if cr.config.Sync {
		sw, err := newSyncWatcher()
		if err != nil {
			return fmt.Errorf("sync watcher: %w", err)
		}
		defer sw.close()
		go sw.watch(ctx, reload)
	}

	if cr.config.Interactive {
		go readStdinReload(ctx, reload)
		log.Info("[chimp] press Enter to reload")
	}

	// Initial run
	cr.execute()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-reload:
			drainChannel(reload)
			log.Info("[chimp] change detected, running command...")
			cr.execute()
		}
	}
}

func (cr *CommandRunner) execute() {
	if err := runHook(cr.config.Dir, cr.config.Command); err != nil {
		log.Warnf("[chimp] command failed: %v", err)
		return
	}

	for _, h := range cr.config.PostHooks {
		if err := runHook(cr.config.Dir, h); err != nil {
			log.Warnf("[chimp] post-hook failed: %v", err)
		}
	}

	if cr.config.Master {
		TouchSyncFile()
	}
}
