package timebuf

import (
	"fmt"
	"testing"
	"time"
)

func TestBufer(t *testing.T) {
	b := NewTimeBuffer(time.Duration(0.2 * 1000000000))
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
			time.Sleep(time.Duration(0.01 * 1000000000))
			b.After()
		}
	}()
	<-c
	d := time.Since(start)
	fmt.Println("duration:", d)
}
