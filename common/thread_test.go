package common

import (
	"testing"
	"time"
)

func TestGoFunc(t *testing.T) {
	GoFunc(func() {
		panic("go func execute error")
	})
	t.Log("func:routine_id:", RoutineId())
	go func() {
		t.Log("go func:routine_id:", RoutineId())
	}()

	time.Sleep(time.Millisecond * 20)
}
