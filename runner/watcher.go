package runner

import (
	"bufio"
	"context"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	log "github.com/s00500/env_logger"
)

type watcher struct {
	fsw       *fsnotify.Watcher
	exts      map[string]bool
	ignore    map[string]bool
	debounce  time.Duration
	watchRoot string
}

func newWatcher(dir string, exts []string, ignoreDirs []string, debounce time.Duration) (*watcher, error) {
	fsw, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	extMap := make(map[string]bool, len(exts))
	for _, e := range exts {
		if !strings.HasPrefix(e, ".") {
			e = "." + e
		}
		extMap[e] = true
	}

	ignoreMap := make(map[string]bool, len(ignoreDirs))
	for _, d := range ignoreDirs {
		ignoreMap[d] = true
	}

	w := &watcher{
		fsw:       fsw,
		exts:      extMap,
		ignore:    ignoreMap,
		debounce:  debounce,
		watchRoot: dir,
	}

	if err := w.addRecursive(dir); err != nil {
		fsw.Close()
		return nil, err
	}

	return w, nil
}

func (w *watcher) addRecursive(root string) error {
	return filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil // skip inaccessible dirs
		}
		if !d.IsDir() {
			return nil
		}
		name := d.Name()
		if w.shouldIgnoreDir(name) {
			return filepath.SkipDir
		}
		return w.fsw.Add(path)
	})
}

func (w *watcher) shouldIgnoreDir(name string) bool {
	if w.ignore[name] {
		return true
	}
	if strings.HasPrefix(name, ".") || strings.HasPrefix(name, "_") {
		return true
	}
	return false
}

func (w *watcher) shouldProcess(name string) bool {
	//// Ignore generated templ files
	//if strings.HasSuffix(name, "_templ.go") {
	//	return false
	//}
	ext := filepath.Ext(name)
	return w.exts[ext]
}

func (w *watcher) watch(ctx context.Context, reload chan<- struct{}) {
	timer := time.NewTimer(0)
	if !timer.Stop() {
		<-timer.C
	}
	pending := false

	for {
		select {
		case <-ctx.Done():
			return

		case event, ok := <-w.fsw.Events:
			if !ok {
				return
			}

			// Watch newly created directories
			if event.Has(fsnotify.Create) {
				if info, err := os.Stat(event.Name); err == nil && info.IsDir() {
					name := filepath.Base(event.Name)
					if !w.shouldIgnoreDir(name) {
						_ = w.fsw.Add(event.Name)
					}
					continue
				}
			}

			if !w.shouldProcess(filepath.Base(event.Name)) {
				continue
			}

			// Reset debounce timer
			if !pending {
				pending = true
				timer.Reset(w.debounce)
			} else {
				if !timer.Stop() {
					select {
					case <-timer.C:
					default:
					}
				}
				timer.Reset(w.debounce)
			}

		case err, ok := <-w.fsw.Errors:
			if !ok {
				return
			}
			log.Warnf("[chimp] watcher error: %v", err)

		case <-timer.C:
			pending = false
			select {
			case reload <- struct{}{}:
			default:
			}
		}
	}
}

func (w *watcher) close() {
	w.fsw.Close()
}

// readStdinReload reads lines from stdin and sends reload signals on Enter.
func readStdinReload(ctx context.Context, reload chan<- struct{}) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		select {
		case <-ctx.Done():
			return
		case reload <- struct{}{}:
		default:
		}
	}
}
