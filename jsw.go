package main

import (
	"./jekyll"
	"./watcher"
	"path/filepath"
	"log"
	"strings"
)

func shouldIgnore(path string) (r bool) {
	s := filepath.ToSlash(path)
	if strings.HasPrefix(s, "_site/") {
		r = true
	}
	return
}

func main() {
	// start jekyll serve
	j := jekyll.NewJekyll()
	err := j.Start()
	if err != nil {
		panic(err)
	}
	// start watcher
	w, err := watcher.NewWatcher(".", nil)
	if err != nil {
		panic(err)
	}
	// infinite loop
	for {
		select {
		case path := <-w.Path:
			log.Println("path:", path)
			if shouldIgnore(path) {
				continue
			}
			// TODO: j.Build()
		}
	}
}
