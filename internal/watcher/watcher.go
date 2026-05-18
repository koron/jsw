// Package watcher wraps fswatcher with recursive directory watching and an
// optional exclude filter. Delivers relative paths on Path and errors on
// Error.
package watcher

import (
	"path/filepath"

	"github.com/fswatcher/fswatcher"
)

// ExcludeFunc is a predicate that returns true for paths that should be
// ignored.
type ExcludeFunc func(path string) bool

// Watcher monitors a directory tree for file system events.
type Watcher struct {
	watcher *fswatcher.Watcher
	Path    chan string
	Error   chan error
}

// noExclude is the default exclude func that allows all paths.
func noExclude(string) bool {
	return false
}

// NewWatcher creates a Watcher for root, recursively watching all
// subdirectories. Events are normalized to paths relative to root; the
// exclude predicate filters out unwanted paths before they reach Path.
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
