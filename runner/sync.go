package runner

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/fsnotify/fsnotify"
	log "github.com/s00500/env_logger"
)

const syncFile = "/tmp/chimp-sync"

// touchSyncFile writes the current timestamp to the sync file (master/producer).
func touchSyncFile() {
	err := os.WriteFile(syncFile, []byte(fmt.Sprintf("%d", time.Now().UnixNano())), 0644)
	if err != nil {
		log.Warnf("[chimp] failed to touch sync file: %v", err)
	}
}

// syncWatcher watches the sync file for changes (consumer).
type syncWatcher struct {
	fsw *fsnotify.Watcher
}

func newSyncWatcher() (*syncWatcher, error) {
	// Ensure the sync file exists so we can watch the parent dir for changes to it
	if _, err := os.Stat(syncFile); os.IsNotExist(err) {
		_ = os.WriteFile(syncFile, []byte("0"), 0644)
	}

	fsw, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	// Watch /tmp and filter to our specific file
	if err := fsw.Add("/tmp"); err != nil {
		fsw.Close()
		return nil, err
	}

	return &syncWatcher{fsw: fsw}, nil
}

func (sw *syncWatcher) watch(ctx context.Context, reload chan<- struct{}) {
	for {
		select {
		case <-ctx.Done():
			return
		case event, ok := <-sw.fsw.Events:
			if !ok {
				return
			}
			if event.Name != syncFile {
				continue
			}
			if event.Has(fsnotify.Write) || event.Has(fsnotify.Create) {
				log.Info("[chimp] sync file changed, reloading...")
				select {
				case reload <- struct{}{}:
				default:
				}
			}
		case err, ok := <-sw.fsw.Errors:
			if !ok {
				return
			}
			log.Warnf("[chimp] sync watcher error: %v", err)
		}
	}
}

func (sw *syncWatcher) close() {
	sw.fsw.Close()
}
