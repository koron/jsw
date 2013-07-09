package main

import (
	"./watcher"
	"log"
)

func main() {
	//ch := make(chan bool)
	w, err := watcher.NewWatcher(".", nil)
	if err != nil {
		panic(err)
	}
	for {
		select {
		case path := <-w.Path:
			log.Println("path:", path)
		}
	}
	//<-ch
}
