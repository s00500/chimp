//go:build !windows

package runner

import (
	"os"
	"syscall"
)

func notifySignals() []os.Signal {
	return []os.Signal{os.Interrupt, syscall.SIGTERM}
}
