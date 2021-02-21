package limit

import (
	"github.com/tiptok/gocomm/sync/task"
	"sync/atomic"
	"testing"
	"time"
)

func TestSlidingWidowLimiter_Allow(t *testing.T) {
	limit, _ := NewSlidingWidowLimiter(time.Second, 10, func() (Window, StopFunc) {
		return NewLocalWindow()
	})
	gt := task.NewGroupTask()
	gt.WithWorkerNumber(2)
	var allowNum int32
	workFunc := func() {
		for i := 0; i < 10; i++ {
			ok := limit.Allow() //"test_counter_limit"
			if ok {
				atomic.AddInt32(&allowNum, 1)
			}
		}
	}
	gt.Run(workFunc)
	gt.Run(workFunc)
	gt.Wait()

	if allowNum != 10 {
		t.Fatalf("allowNum want:%d get:%d", 10, allowNum)
	}
}
