package task

import (
	"sync/atomic"
	"testing"
	"time"
)

func TestGroupTask_Run(t *testing.T) {
	group := NewGroupTask()
	group.WithWorkerNumber(4)
	var sum int32 = 0
	for i := 0; i < 12; i++ {
		group.Run(func() {
			time.Sleep(time.Millisecond * 50)
			atomic.AddInt32(&sum, 1)
		})
	}
	group.Wait()
	if sum != 12 {
		t.Fatal("expect:", 3, "get:", sum)
	}
}
