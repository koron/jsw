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

func regulatePath(path string) (r string) {
	r = filepath.ToSlash(path)
	if strings.HasPrefix(r, "./") {
		r = r[2:]
	}
	return
}

func shouldIgnore(path string) (r bool) {
	s := regulatePath(path)
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
	m := make(map[string]int)
	for {
		select {
		case path := <-w.Path:
			path = regulatePath(path)
			_, present := m[path]
			if !present {
				m[path] = 1
				log.Println("changed:", path)
			}
			tb.After()
		case _ = <-tb.C:
			log.Println("build started: changed", len(m), "files")
			j.Build()
			log.Println("build finished")
			m = make(map[string]int)
		case _ = <-w.Error:
		}
	}
}
