package main

import (
	"github.com/howeyc/fsnotify"
	"io"
	"log"
	"os"
	"os/exec"
)

func startServer() *exec.Cmd {
	// Start jekyll server.
	cmd := exec.Command("jekyll", "serve")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		panic(err)
	}
	err = cmd.Start()
	if err != nil {
		panic(err)
	}
	go io.Copy(os.Stdout, stdout)
	go io.Copy(os.Stderr, stderr)
	return cmd
}

func main() {
	serv_cmd := startServer()
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		panic(err)
	}

	err = watcher.Watch(".")
	if err != nil {
		panic(err)
	}

	for {
		select {
		case ev := <-watcher.Event:
			log.Println("event:", ev)
		case err := <-watcher.Error:
			panic(err)
		}
	}
	serv_cmd.Wait()
}
