//go:build !windows

package runner

import (
	"os/exec"
	"syscall"
)

func setProcAttr(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
}

func gracefulStop(cmd *exec.Cmd) {
	// Send SIGINT to the process group
	_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGINT)
}

func forceStop(cmd *exec.Cmd) {
	_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
}
