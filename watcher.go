package main

import (
	"github.com/howeyc/fsnotify"
	"log"
	"os"
	"path/filepath"
)

func startWatcher(root string) (w *fsnotify.Watcher, err error) {
	w, err = fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	addDir := func(path string) {
		w.Watch(path)
	}
	removeDir := func(path string) {
		w.RemoveWatch(path)
	}

	// register initial directories to watch.
	err = filepath.Walk(root,
		func(path string, f os.FileInfo, err error) error {
			if f.IsDir() {
				addDir(path)
			}
			return err
		})
	if err != nil {
		w.Close()
		return nil, err
	}

	emitEvent := func(ev *fsnotify.FileEvent) {
		log.Println(ev)
	}

	emitError := func(err error) {
		panic(err)
	}

	go func() {
		for {
			select {
			case ev := <-w.Event:
				name := ev.Name
				if ev.IsCreate() {
					fi, err := os.Stat(name)
					if err != nil {
						emitError(err)
						break
					} else if fi.IsDir() {
						addDir(name)
					}
				} else if ev.IsDelete() {
					removeDir(name)
				}
				emitEvent(ev)
			case err := <-w.Error:
				emitError(err)
			}
		}
	}()
	return
}

func main() {
	ch := make(chan bool)
	startWatcher(".")
	<-ch
}
