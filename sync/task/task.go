package task

import (
	"context"

	"github.com/tiptok/gocomm/sync/signal/semaphore"
)

// OnSuccess executes g() after f() returns nil.
// 执行两个任务 一个前置 beforeWork 成功以后执行 work
func OnSuccess(beforeWork func() error, work func() error) func() error {
	return func() error {
		if err := beforeWork(); err != nil {
			return err
		}
		return work()
	}
}

// Run executes a list of tasks in parallel, returns the first error encountered or nil if all tasks pass.
// 并行执行一串任务
func Run(ctx context.Context, tasks ...func() error) error {
	n := len(tasks)
	s := semaphore.New(n)
	done := make(chan error, 1)

	for i := range tasks {
		<-s.Wait()
		go func(f func() error) {
			err := f()
			if err == nil {
				s.Signal()
				return
			}

			select {
			case done <- err:
			default:
			}
		}(tasks[i])
	}

	for i := 0; i < n; i++ {
		select {
		case err := <-done:
			return err
		case <-ctx.Done():
			return ctx.Err()
		case <-s.Wait():
		}
	}

	return nil
}
