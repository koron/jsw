// Command jsw starts `jekyll serve` then watches the project directory for
// file changes and triggers `jekyll build` on a debounced schedule.
package main

import (
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"time"

	"github.com/koron/jsw/internal/jekyll"
	"github.com/koron/jsw/internal/timebuf"
	"github.com/koron/jsw/internal/watcher"
)

// regulatePath normalizes a path: converts backslashes to slashes and strips
// a leading "./".
func regulatePath(path string) (r string) {
	r = filepath.ToSlash(path)
	if strings.HasPrefix(r, "./") {
		r = r[2:]
	}
	return
}

// shouldIgnore returns true for paths under _site/ or .git/.
func shouldIgnore(path string) (r bool) {
	s := regulatePath(path)
	if strings.HasPrefix(s, "_site") || strings.HasPrefix(s, ".git") {
		r = true
	}
	return
}

// main starts jekyll serve, sets up a recursive file watcher with a 200 ms
// debounce, and loops on change events to call jekyll build. SIGINT
// gracefully stops the server.
func main() {
	// start jekyll serve
	j := jekyll.NewJekyll()
	err := j.Start()
	if err != nil {
		panic(err)
	}

	// handle SIGINT
	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, os.Interrupt)
		<-sig
		j.Stop()
		os.Exit(0)
	}()

	// start watcher
	w, err := watcher.NewWatcher(".", shouldIgnore)
	if err != nil {
		panic(err)
	}

	// prepare a timer to build
	tb := timebuf.NewTimeBuffer(time.Duration(200 * time.Millisecond))

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
