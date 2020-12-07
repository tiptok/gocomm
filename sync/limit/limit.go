package limit

import "errors"

var ErrReturn = errors.New("discarding limited token, resource pool is full, someone returned multiple times")

type Limit struct {
	pool chan interface{}
}

func NewLimit(n int) Limit {
	return Limit{
		pool: make(chan interface{}, n),
	}
}

func (l Limit) Borrow() {
	l.pool <- struct{}{}
}

// Return returns the borrowed resource, returns error only if returned more than borrowed.
func (l Limit) Return() error {
	select {
	case <-l.pool:
		return nil
	default:
		return ErrReturn
	}
}

func (l Limit) TryBorrow() bool {
	select {
	case l.pool <- struct{}{}:
		return true
	default:
		return false
	}
}
