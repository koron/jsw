package watcher

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func recvEventOrTimeout(ch <-chan string, d time.Duration) (string, bool) {
	select {
	case v := <-ch:
		return v, true
	case <-time.After(d):
		return "", false
	}
}

func recvErrorOrTimeout(ch <-chan error, d time.Duration) (error, bool) {
	select {
	case v := <-ch:
		return v, true
	case <-time.After(d):
		return nil, false
	}
}

// TestCreateFile verifies that creating a file under the watched directory
// produces an event on the Path channel.
func TestCreateFile(t *testing.T) {
	dir := t.TempDir()

	w, err := NewWatcher(dir, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Give the watcher time to set up OS-level watches.
	time.Sleep(500 * time.Millisecond)

	f, err := os.Create(filepath.Join(dir, "test.txt"))
	if err != nil {
		t.Fatal(err)
	}
	f.Close()

	got, ok := recvEventOrTimeout(w.Path, 3*time.Second)
	if !ok {
		t.Fatal("expected file creation event, got none")
	}
	if !strings.HasSuffix(got, "test.txt") {
		t.Fatalf("expected path ending with test.txt, got %q", got)
	}
}

// TestExcludedPath verifies that events for paths rejected by the exclude
// function are not delivered on the Path channel.
func TestExcludedPath(t *testing.T) {
	dir := t.TempDir()

	exclude := func(path string) bool {
		return strings.HasPrefix(path, "_excluded")
	}
	w, err := NewWatcher(dir, exclude)
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(500 * time.Millisecond)

	// Create a file under an excluded directory.
	exDir := filepath.Join(dir, "_excluded")
	if err := os.MkdirAll(exDir, 0755); err != nil {
		t.Fatal(err)
	}

	f, err := os.Create(filepath.Join(exDir, "ignored.txt"))
	if err != nil {
		t.Fatal(err)
	}
	f.Close()

	// The excluded event must not arrive on Path.
	if ev, ok := recvEventOrTimeout(w.Path, 3*time.Second); ok {
		t.Fatalf("expected no event for excluded path, got %q", ev)
	}
}

// TestCreateFileInSubdir verifies that a file created in a subdirectory
// is detected by the recursive watcher.
func TestCreateFileInSubdir(t *testing.T) {
	dir := t.TempDir()

	w, err := NewWatcher(dir, nil)
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(500 * time.Millisecond)

	sub := filepath.Join(dir, "subdir")
	if err := os.MkdirAll(sub, 0755); err != nil {
		t.Fatal(err)
	}
	// Creating the subdirectory may trigger an event; drain it.
	recvEventOrTimeout(w.Path, 3*time.Second)

	f, err := os.Create(filepath.Join(sub, "nested.txt"))
	if err != nil {
		t.Fatal(err)
	}
	f.Close()

	got, ok := recvEventOrTimeout(w.Path, 3*time.Second)
	if !ok {
		t.Fatal("expected event for file in subdirectory, got none")
	}
	if !strings.HasSuffix(got, "nested.txt") {
		t.Fatalf("expected path ending with nested.txt, got %q", got)
	}
}
