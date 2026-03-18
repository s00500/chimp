//go:build windows

package runner

import "os"

func notifySignals() []os.Signal {
	return []os.Signal{os.Interrupt}
}
