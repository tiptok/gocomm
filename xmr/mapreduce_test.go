package xmr

import (
	"errors"
	"github.com/tiptok/gocomm/common"
	"runtime"
	"sync/atomic"
	"testing"
	"time"
)

var errDummy = errors.New("dummy")

func TestMapReduce(t *testing.T) {
	tests := []struct {
		mapper      MapperFunc
		reducer     ReducerFunc
		expectErr   error
		expectValue interface{}
	}{
		{
			expectErr:   nil,
			expectValue: 30,
		},
		{
			mapper: func(item interface{}, writer Writer, cancel func(error)) {
				v := item.(int)
				if v%3 == 0 {
					cancel(errDummy)
				}
				writer.Write(v * v)
			},
			expectErr: errDummy,
		},
		{
			mapper: func(item interface{}, writer Writer, cancel func(error)) {
				v := item.(int)
				if v%3 == 0 {
					cancel(nil)
				}
				writer.Write(v * v)
			},
			expectErr:   ErrCancelWithNil,
			expectValue: nil,
		},
		{
			reducer: func(pipe <-chan interface{}, writer Writer, cancel func(error)) {
				var result int
				for item := range pipe {
					result += item.(int)
					if result > 10 {
						cancel(errDummy)
					}
				}
				writer.Write(result)
			},
			expectErr: errDummy,
		},
	}

	for _, test := range tests {
		t.Run(common.RandomString(8), func(t *testing.T) {
			if test.mapper == nil {
				test.mapper = func(item interface{}, writer Writer, cancel func(error)) {
					v := item.(int)
					writer.Write(v * v)
				}
			}
			if test.reducer == nil {
				test.reducer = func(pipe <-chan interface{}, writer Writer, cancel func(error)) {
					var result int
					for item := range pipe {
						result += item.(int)
					}
					writer.Write(result)
				}
			}
			value, err := MapReduce(func(source chan<- interface{}) {
				for i := 1; i < 5; i++ {
					source <- i
				}
			}, test.mapper, test.reducer, WithWorkers(runtime.NumCPU()))

			if err != test.expectErr {
				t.Fatal(err, " except:", test.expectErr)
			}
			if value != test.expectValue {
				t.Fatal(err, " except:", test.expectErr)
			}
			//assert.Equal(t, test.expectErr, err)
			//assert.Equal(t, test.expectValue, value)
		})
	}
}

func TestExampleReduce(t *testing.T) {
	var value int32
	add := func() error {
		time.Sleep(time.Millisecond * 100)
		atomic.AddInt32(&value, 1)
		return nil
	}
	err := Reduce(add, add, add, add)
	if err != nil {
		t.Fatal(err)
	}
	if value != 4 {
		t.Fatal("except:value=4 get:", value)
	}
}
