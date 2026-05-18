package timebuf

import (
	"time"
)

type TimeBuffer struct {
	C  chan time.Time
	d  time.Duration
	id int
}

func NewTimeBuffer(d time.Duration) (b *TimeBuffer) {
	b = &TimeBuffer{
		C: make(chan time.Time, 0),
		d: d,
	}
	return
}

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
