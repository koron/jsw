package watcher

import (
	"github.com/howeyc/fsnotify"
	"os"
	"path/filepath"
)

type ExcludeFunc func(path string) bool

type Watcher struct {
	watcher *fsnotify.Watcher
	Path    chan string
	Error   chan error
}

func noExclude(string) bool {
	return false
}

func NewWatcher(root string, exclude ExcludeFunc) (w *Watcher, err error) {
	fsw, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	if exclude == nil {
		exclude = noExclude
	}

	addDir := func(path string) {
		fsw.Watch(path)
	}

	removeDir := func(path string) {
		fsw.RemoveWatch(path)
	}

	// register initial directories to watch.
	err = filepath.Walk(root,
		func(path string, f os.FileInfo, err error) error {
			if f.IsDir() && !exclude(path) {
				addDir(path)
			}
			return err
		})
	if err != nil {
		fsw.Close()
		return nil, err
	}

	w = &Watcher{fsw, make(chan string), make(chan error)}

	go func() {
		for {
			select {
			case ev := <-fsw.Event:
				name := ev.Name
				if ev.IsCreate() {
					fi, err := os.Stat(name)
					if err != nil {
						w.Error <- err
						break
					} else if fi.IsDir() {
						addDir(name)
					}
				} else if ev.IsDelete() {
					removeDir(name)
				}

				if !exclude(name) {
					w.Path <- name
				}
			case err := <-fsw.Error:
				w.Error <- err
			}
		}
	}()

	return
}
