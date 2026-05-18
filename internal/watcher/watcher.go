package watcher

import (
	"path/filepath"

	"github.com/fswatcher/fswatcher"
)

type ExcludeFunc func(path string) bool

type Watcher struct {
	watcher *fswatcher.Watcher
	Path    chan string
	Error   chan error
}

func noExclude(string) bool {
	return false
}

func NewWatcher(root string, exclude ExcludeFunc) (w *Watcher, err error) {
	if exclude == nil {
		exclude = noExclude
	}
	rootabs, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}

	watcher, err := fswatcher.NewWatcher()
	if err != nil {
		return nil, err
	}
	watcher.AddRecursive(rootabs, fswatcher.All)

	w = &Watcher{watcher, make(chan string), make(chan error)}

	go func() {
		for {
			select {
			case ev := <-watcher.Events:
				name, err := filepath.Rel(rootabs, ev.Name)
				if err != nil {
					w.Error <- err
					break
				}
				if !exclude(name) {
					w.Path <- name
				}
			case err := <-watcher.Errors:
				w.Error <- err
			}
		}
		watcher.Close()
	}()

	return w, nil
}
