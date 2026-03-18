//go:build windows

package runner

import (
	"fmt"
	"os/exec"
)

func setProcAttr(cmd *exec.Cmd) {
	// No process group attributes needed on Windows
}

func gracefulStop(cmd *exec.Cmd) {
	// On Windows, there is no reliable way to send a graceful signal to a
	// console process tree. Use taskkill /T to kill the process tree.
	_ = exec.Command("taskkill", "/F", "/T", "/PID", fmt.Sprintf("%d", cmd.Process.Pid)).Run()
}

func forceStop(cmd *exec.Cmd) {
	_ = cmd.Process.Kill()
}
