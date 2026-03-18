package runner

import (
	"context"
	"os"
	"os/exec"
	"time"

	log "github.com/s00500/env_logger"
)

type process struct {
	dir         string
	args        []string
	interactive bool
	cmd         *exec.Cmd
	done        chan struct{}
}

func newProcess(dir string, args []string, interactive bool) *process {
	return &process{
		dir:         dir,
		args:        args,
		interactive: interactive,
	}
}

func (p *process) start(ctx context.Context) {
	cmdArgs := append([]string{"run", "."}, p.args...)
	p.cmd = exec.CommandContext(ctx, "go", cmdArgs...)
	p.cmd.Dir = p.dir
	p.cmd.Stdout = os.Stdout
	p.cmd.Stderr = os.Stderr
	setProcAttr(p.cmd)

	if !p.interactive {
		p.cmd.Stdin = os.Stdin
	}

	p.done = make(chan struct{})

	if err := p.cmd.Start(); err != nil {
		log.Errorf("[chimp] failed to start: %v", err)
		close(p.done)
		return
	}

	log.Infof("[chimp] started (pid %d)", p.cmd.Process.Pid)

	go func() {
		defer close(p.done)
		err := p.cmd.Wait()
		if err != nil {
			log.Warnf("[chimp] process exited: %v", err)
		} else {
			log.Info("[chimp] process exited successfully")
		}
	}()
}

func (p *process) stop() {
	if p.cmd == nil || p.cmd.Process == nil {
		return
	}

	gracefulStop(p.cmd)

	// Wait up to 5 seconds for graceful shutdown
	select {
	case <-p.done:
		return
	case <-time.After(5 * time.Second):
		log.Warn("[chimp] process did not exit gracefully, force killing")
		forceStop(p.cmd)
		<-p.done
	}
}

// runHook executes a single hook command.
func runHook(dir string, h HookConfig) error {
	cmd := exec.Command(h.Command, h.Args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
