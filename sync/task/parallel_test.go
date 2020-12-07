package task

import (
	"sync/atomic"
	"testing"
	"time"
)

func TestParallel(t *testing.T) {
	var sum int32 = 0
	fn := func() {
		time.Sleep(time.Millisecond * 50)
		atomic.AddInt32(&sum, 1)
	}
	Parallel(fn, fn, fn, fn)
	if sum != 4 {
		t.Fatal("expect:", 4, "get:", sum)
	}
}
