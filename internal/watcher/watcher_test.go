package watcher

import (
	"log"
	"testing"
)

func TestWatcher(t *testing.T) {
	w, err := NewWatcher(".", nil)
	if err != nil {
		t.Error(err)
	}
	for {
		select {
		case path := <-w.Path:
			log.Println(path)
		case err = <-w.Error:
			log.Println(err)
		}
	}
}
