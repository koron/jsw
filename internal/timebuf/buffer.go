// Package timebuf provides a debounce timer that coalesces rapid
// consecutive events into a single notification after a quiet period.
package timebuf

import (
	"time"
)

// TimeBuffer is a debounce timer. Each call to After resets the window;
// the channel C receives a value once no new calls arrive within the
// configured duration.
type TimeBuffer struct {
	C  chan time.Time
	d  time.Duration
	id int
}

// NewTimeBuffer creates a TimeBuffer with the given debounce duration d.
// The returned channel C delivers one time.Time per debounced batch.
func NewTimeBuffer(d time.Duration) (b *TimeBuffer) {
	b = &TimeBuffer{
		C: make(chan time.Time, 0),
		d: d,
	}
	return
}

// After resets the debounce window. If no other call to After is made
// within the configured duration, a timestamp is sent on C.
func (b *TimeBuffer) After() {
	b.id++
	id := b.id
	go func(target int) {
		time.Sleep(b.d)
		if target == b.id {
			b.C <- time.Now()
		}
	}(id)
}
