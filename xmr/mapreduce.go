package xmr

import (
	"errors"
	"fmt"
	"github.com/tiptok/gocomm/common"
	"github.com/tiptok/gocomm/xerror"
	"sync"
)

const (
	defaultWorkers = 16
	minWorkers     = 1
)

var (
	ErrCancelWithNil  = errors.New("mapreduce cancelled with nil")
	ErrReduceNoOutput = errors.New("reduce not writing value")
)

type (
	GenerateFunc func(source chan<- interface{})
	MapFunc      func(item interface{}, writer Writer)
	//VoidMapFunc     func(item interface{})
	MapperFunc  func(item interface{}, writer Writer, cancel func(error))
	ReducerFunc func(pipe <-chan interface{}, writer Writer, cancel func(error))
	//VoidReducerFunc func(pipe <-chan interface{}, cancel func(error))
	Option func(opts *mapReduceOptions)

	mapReduceOptions struct {
		workers int
	}

	Writer interface {
		Write(v interface{})
	}
)

func Reduce(fns ...func() error) error {
	if len(fns) == 0 {
		return nil
	}

	_, err := MapReduce(func(source chan<- interface{}) {
		for _, fn := range fns {
			source <- fn
		}
	}, func(item interface{}, writer Writer, cancel func(error)) {
		fn := item.(func() error)
		if err := fn(); err != nil {
			cancel(err)
		}
	}, func(pipe <-chan interface{}, writer Writer, cancel func(error)) {
		drain(pipe)
		writer.Write(struct {
		}{})
	}, WithWorkers(len(fns)))

	return err
}

func MapReduce(generate GenerateFunc, mapper MapperFunc, reducer ReducerFunc, opts ...Option) (interface{}, error) {
	source := buildSource(generate)
	return MapReduceWithSource(source, mapper, reducer, opts...)
}

func buildSource(generate GenerateFunc) chan interface{} {
	source := make(chan interface{})
	common.GoFunc(func() {
		defer close(source)
		generate(source)
	})
	return source
}

func MapReduceWithSource(source <-chan interface{}, mapper MapperFunc, reducer ReducerFunc,
	opts ...Option) (interface{}, error) {
	options := buildOptions(opts...)
	output := make(chan interface{})
	collector := make(chan interface{}, options.workers)
	done := make(chan struct{})
	writer := newGuardedWriter(output, done)
	var closeOnce sync.Once
	var retErr xerror.AtomicError
	finish := func() {
		closeOnce.Do(func() {
			close(done)
			close(output)
		})
	}
	cancel := once(func(err error) {
		if err != nil {
			retErr.Set(err)
		} else {
			retErr.Set(ErrCancelWithNil)
		}

		drain(source)
		finish()
	})

	go func() {
		defer func() {
			if r := recover(); r != nil {
				cancel(fmt.Errorf("%v", r))
			} else {
				finish()
			}
		}()
		reducer(collector, writer, cancel)
		drain(collector)
	}()

	go executeMappers(func(item interface{}, w Writer) {
		mapper(item, w, cancel)
	}, source, collector, done, options.workers)

	value, ok := <-output
	if err := retErr.Load(); err != nil {
		return nil, err
	} else if ok {
		return value, nil
	} else {
		return nil, ErrReduceNoOutput
	}
}

func executeMappers(mapper MapFunc, input <-chan interface{}, collector chan<- interface{},
	done <-chan struct{}, workers int) {
	var wg sync.WaitGroup
	defer func() {
		wg.Wait()
		close(collector)
	}()

	pool := make(chan struct{}, workers)
	writer := newGuardedWriter(collector, done)
	for {
		select {
		case <-done:
			return
		case pool <- struct{}{}:
			item, ok := <-input
			if !ok {
				<-pool
				return
			}

			wg.Add(1)
			// better to safely run caller defined method
			common.GoFunc(func() {
				defer func() {
					wg.Done()
					<-pool
				}()

				mapper(item, writer)
			})
		}
	}
}

func once(fn func(error)) func(error) {
	once := new(sync.Once)
	return func(err error) {
		once.Do(func() {
			fn(err)
		})
	}
}

// drain drains the channel.
func drain(channel <-chan interface{}) {
	// drain the channel
	for range channel {
	}
}

func newOptions() *mapReduceOptions {
	return &mapReduceOptions{
		workers: defaultWorkers,
	}
}

func WithWorkers(workers int) Option {
	return func(opts *mapReduceOptions) {
		if workers < minWorkers {
			opts.workers = minWorkers
		} else {
			opts.workers = workers
		}
	}
}

func buildOptions(opts ...Option) *mapReduceOptions {
	options := newOptions()
	for _, opt := range opts {
		opt(options)
	}

	return options
}

type guardedWriter struct {
	channel chan<- interface{}
	done    <-chan struct{}
}

func newGuardedWriter(channel chan<- interface{}, done <-chan struct{}) guardedWriter {
	return guardedWriter{
		channel: channel,
		done:    done,
	}
}
func (gw guardedWriter) Write(v interface{}) {
	select {
	case <-gw.done:
		return
	default:
		gw.channel <- v
	}
}
