package limit

import (
	"github.com/tiptok/gocomm/sync/task"
	"sync/atomic"
	"testing"
	"time"
)

func TestNewCounterLimit(t *testing.T) {
	limit := NewCounterLimiter(10, time.Second)
	gt := task.NewGroupTask()
	gt.WithWorkerNumber(2)
	//var gNum int32
	var allowNum int32
	workFunc := func() {
		//num:= atomic.AddInt32(&gNum,1)
		for i := 0; i < 10; i++ {
			allow := limit.Allow("test_counter_limit")
			if allow {
				//t.Logf("groutine(%d) take:%s sn:%d allowed",num,"test_counter_limit",i)
				atomic.AddInt32(&allowNum, 1)
			} else {
				//t.Logf("groutine(%d) take:%s sn:%d OverQuota",num,"test_counter_limit",i)
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
