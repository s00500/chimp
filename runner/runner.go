package runner

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/s00500/env_logger"
)

type Config struct {
	Args        []string      // passed to the compiled binary as args
	Dir         string        // working dir (default: cwd)
	WatchExts   []string      // default: [".go"]
	IgnoreDirs  []string      // default: .git, vendor, node_modules, tmp
	PreHooks    []HookConfig  // run before each restart
	PostHooks   []HookConfig  // run after each restart
	Master      bool          // touch the sync file on reload (producer)
	Sync        bool          // watch the sync file for reload signals (consumer)
	Interactive bool          // manual reload key (default: true)
	Debounce    time.Duration // default: 200ms
}

type HookConfig struct {
	Command string
	Args    []string
}

type Runner struct {
	config  Config
	process *process
}

func New(config Config) *Runner {
	if config.Dir == "" {
		config.Dir, _ = os.Getwd()
	}
	if len(config.WatchExts) == 0 {
		config.WatchExts = []string{".go"}
	}
	if len(config.IgnoreDirs) == 0 {
		config.IgnoreDirs = []string{".git", "vendor", "node_modules", "tmp"}
	}
	if config.Debounce == 0 {
		config.Debounce = 200 * time.Millisecond
	}
	return &Runner{config: config}
}

func (r *Runner) Run(ctx context.Context) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer cancel()

	reload := make(chan struct{}, 1)

	// Start file watcher
	w, err := newWatcher(r.config.Dir, r.config.WatchExts, r.config.IgnoreDirs, r.config.Debounce)
	if err != nil {
		return fmt.Errorf("watcher: %w", err)
	}
	defer w.close()

	go w.watch(ctx, reload)

	// Start sync watcher if consumer
	if r.config.Sync {
		sw, err := newSyncWatcher()
		if err != nil {
			return fmt.Errorf("sync watcher: %w", err)
		}
		defer sw.close()
		go sw.watch(ctx, reload)
	}

	// Interactive reload via Enter key
	if r.config.Interactive {
		go readStdinReload(ctx, reload)
		log.Info("[chimp] press Enter to reload")
	}

	// Initial build and run
	r.restart(ctx)

	for {
		select {
		case <-ctx.Done():
			r.stop()
			return nil
		case <-reload:
			drainChannel(reload)
			log.Info("[chimp] reloading...")
			r.stop()
			r.runHooks(r.config.PreHooks, "pre-hook")
			r.start(ctx)
			r.runHooks(r.config.PostHooks, "post-hook")

			if r.config.Master {
				touchSyncFile()
			}
		}
	}
}

func (r *Runner) restart(ctx context.Context) {
	r.stop()
	r.runHooks(r.config.PreHooks, "pre-hook")
	r.start(ctx)
	r.runHooks(r.config.PostHooks, "post-hook")

	if r.config.Master {
		touchSyncFile()
	}
}

func (r *Runner) start(ctx context.Context) {
	p := newProcess(r.config.Dir, r.config.Args, r.config.Interactive)
	r.process = p
	p.start(ctx)
}

func (r *Runner) stop() {
	if r.process != nil {
		r.process.stop()
		r.process = nil
	}
}

func (r *Runner) runHooks(hooks []HookConfig, label string) {
	for _, h := range hooks {
		log.Infof("[chimp] running %s: %s %v", label, h.Command, h.Args)
		if err := runHook(r.config.Dir, h); err != nil {
			log.Warnf("[chimp] %s failed: %v", label, err)
		}
	}
}

func drainChannel(ch chan struct{}) {
	for {
		select {
		case <-ch:
		default:
			return
		}
	}
}
