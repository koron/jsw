package watcher

import (
	"log"
	"testing"
)

// TestWatcher is a placeholder; skipped because it requires manual
// inspection of live file events.
func TestWatcher(t *testing.T) {
	t.Skip("not implemented yet")
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
