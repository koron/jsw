package main

import (
	"./jekyll"
	"./timebuf"
	"./watcher"
	"log"
	"path/filepath"
	"strings"
	"time"
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
	// prepare a timer to build
	tb := timebuf.NewTimeBuffer(time.Duration(0.2 * 1000000000))

	// infinite loop
	for {
		select {
		case path := <-w.Path:
			//log.Println("path:", path)
			if shouldIgnore(path) {
				continue
			}
			tb.After()
		case _ = <-w.Error:
			//log.Println("warn:", err)
		case _ = <-tb.C:
			log.Println("rebuilding")
			j.Build()
			log.Println("rebuilded")
		}
	}
}
