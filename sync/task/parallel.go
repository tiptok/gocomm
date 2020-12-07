package task

func Parallel(fns ...func()) {
	group := NewGroupTask()
	for _, fn := range fns {
		group.Run(fn)
	}
	group.Wait()
}
