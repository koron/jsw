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
	if strings.HasPrefix(s, "./") {
		s = s[2:]
	}
	if strings.HasPrefix(s, "_site") || strings.HasPrefix(s, ".git") {
		r = true
	}
	//log.Println("be ignored", s, r)
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
	w, err := watcher.NewWatcher(".", shouldIgnore)
	if err != nil {
		panic(err)
	}
	// prepare a timer to build
	tb := timebuf.NewTimeBuffer(time.Duration(0.2 * 1000000000))

	// infinite loop
	for {
		select {
		case _ = <-w.Path:
			tb.After()
		case _ = <-tb.C:
			log.Println("build started")
			j.Build()
			log.Println("build finished")
		case _ = <-w.Error:
		}
	}
}
