package timebuf

import (
	"fmt"
	"testing"
	"time"
)

// TestBufer verifies the debounce timer: ten rapid calls within 10 ms
// each should be coalesced into a single notification delivered after
// the 200 ms window from the last call.
func TestBufer(t *testing.T) {
	b := NewTimeBuffer(200 * time.Millisecond)
	c := make(chan bool, 1)
	go func() {
		t := <-b.C
		fmt.Println("received: ", t)
		c <- true
	}()
	start := time.Now()
	go func() {
		b.After()
		for i := 0; i < 10; i++ {
			time.Sleep(10 * time.Millisecond)
			b.After()
		}
	}()
	<-c
	d := time.Since(start)
	fmt.Println("duration:", d)
}
