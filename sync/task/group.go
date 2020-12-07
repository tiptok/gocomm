package task

import (
	"github.com/tiptok/gocomm/common"
	"sync"
)

type GroupTask struct {
	wg     sync.WaitGroup
	worker *workerLimit
}

func NewGroupTask() *GroupTask {
	return &GroupTask{}
}

// WithWorkerNumber 限制工作线程数，没设置的话有几个执行任务旧创建几个协程
func (g *GroupTask) WithWorkerNumber(num int) {
	g.worker = NewWorkerLimit(num)
}

func (g *GroupTask) Run(fn func()) {
	g.wg.Add(1)
	if g.worker != nil {
		g.worker.Schedule(func() {
			defer g.wg.Done()
			fn()
		})
		return
	}
	common.GoFunc(func() {
		defer g.wg.Done()
		fn()
	})
}

func (g *GroupTask) Wait() {
	g.wg.Wait()
}

type workerLimit struct {
	limitChan chan struct{}
}

func NewWorkerLimit(concurrency int) *workerLimit {
	return &workerLimit{
		limitChan: make(chan struct{}, concurrency),
	}
}

func (wl *workerLimit) Schedule(task func()) {
	wl.limitChan <- struct{}{}

	go func() {
		defer common.Recover(func() {
			<-wl.limitChan
		})

		task()
	}()
}
