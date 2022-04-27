package idgen

import (
	"testing"
)

func Test_Next(t *testing.T) {
	m := make(map[int64]int64)
	num := 1000
	for i := 0; i < num; i++ {
		id := Next()
		if _, ok := m[id]; ok {
			t.Fatal("exists id:", id, len(m))
		}
		t.Log(id, Decompose(uint64(id)))
	}
}

func TestDecompose(t *testing.T) {
	Init(func() (uint16, error) {
		return 2, nil
	})
	id := Next()
	for i := 0; i < 100; i++ {
		id = Next()
	}
	t.Log(id)
	t.Log(Decompose(uint64(id)))
}
