package xbyte

import (
	"github.com/tiptok/gocomm/sync/task"
	"io"
	"os"
	"sync"
	"testing"
)

func BenchmarkNewBytePoolCap(b *testing.B) {
	bytePool := NewBytePoolCap(500, 1024, 1024)
	for i := 0; i < b.N; i++ {
		opBytePool(bytePool)
	}
}

func opBytePool(bytePool *BytePoolCap) {
	t := task.NewGroupTask()
	t.WithWorkerNumber(500)
	t.Run(func() {
		buffer := bytePool.Get()
		defer bytePool.Put(buffer)
		mockReadFile(buffer)
	})
	t.Wait()
}

func BenchmarkNewSyncPoolCap(b *testing.B) {
	bytePool := &sync.Pool{
		New: func() interface{} {
			return make([]byte, 1024, 1024)
		},
	}
	for i := 0; i < b.N; i++ {
		opSyncPool(bytePool)
	}
}

func opSyncPool(bytePool *sync.Pool) {
	t := task.NewGroupTask()
	t.WithWorkerNumber(500)
	t.Run(func() {
		buffer := bytePool.Get()
		defer bytePool.Put(buffer)
		mockReadFile(buffer.([]byte))
	})
	t.Wait()
}

func mockReadFile(b []byte) {
	f, _ := os.Open("F:\\go\\src\\mmm.rar")
	for {
		n, err := io.ReadFull(f, b)
		if n == 0 || err == io.EOF {
			break
		}
	}
}
