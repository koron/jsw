package timebuf

import (
	"testing"
	"time"
)

func recvOrTimeout(t *testing.T, ch <-chan time.Time, d time.Duration) (time.Time, bool) {
	t.Helper()
	select {
	case v := <-ch:
		return v, true
	case <-time.After(d):
		return time.Time{}, false
	}
}

func assertRecv(t *testing.T, ch <-chan time.Time, d time.Duration) time.Time {
	t.Helper()
	v, ok := recvOrTimeout(t, ch, d)
	if !ok {
		t.Fatalf("expected receive within %v, got none", d)
	}
	return v
}

func assertNoRecv(t *testing.T, ch <-chan time.Time, d time.Duration) {
	t.Helper()
	v, ok := recvOrTimeout(t, ch, d)
	if ok {
		t.Fatalf("expected no receive, got %v", v)
	}
}

// TestBufer verifies the debounce timer: ten rapid calls within 10 ms
// each should be coalesced into a single notification delivered after
// the 200 ms window from the last call.
func TestBufer(t *testing.T) {
	b := NewTimeBuffer(200 * time.Millisecond)
	start := time.Now()
	go func() {
		b.After()
		for i := 0; i < 10; i++ {
			time.Sleep(10 * time.Millisecond)
			b.After()
		}
	}()
	assertRecv(t, b.C, 500*time.Millisecond)
	d := time.Since(start)
	if d < 200*time.Millisecond {
		t.Errorf("duration too short: %v", d)
	}
}

// TestSingleCall verifies that a single After() call fires exactly once
// within the debounce duration.
func TestSingleCall(t *testing.T) {
	b := NewTimeBuffer(50 * time.Millisecond)
	start := time.Now()
	b.After()
	assertRecv(t, b.C, 200*time.Millisecond)
	d := time.Since(start)
	if d < 50*time.Millisecond {
		t.Errorf("duration too short: %v", d)
	}
}

// TestNoCall verifies that C never receives when After() is never called.
func TestNoCall(t *testing.T) {
	b := NewTimeBuffer(50 * time.Millisecond)
	assertNoRecv(t, b.C, 200*time.Millisecond)
}

// TestSequentialBatches verifies that two non-overlapping batches produce
// two separate events.
func TestSequentialBatches(t *testing.T) {
	b := NewTimeBuffer(50 * time.Millisecond)

	b.After()
	assertRecv(t, b.C, 200*time.Millisecond)

	time.Sleep(100 * time.Millisecond)

	b.After()
	assertRecv(t, b.C, 200*time.Millisecond)

	assertNoRecv(t, b.C, 200*time.Millisecond)
}

// TestConcurrentCalls verifies that concurrent After() calls are safe.
func TestConcurrentCalls(t *testing.T) {
	b := NewTimeBuffer(50 * time.Millisecond)
	for i := 0; i < 10; i++ {
		go b.After()
	}
	assertRecv(t, b.C, 500*time.Millisecond)
	assertNoRecv(t, b.C, 200*time.Millisecond)
}
