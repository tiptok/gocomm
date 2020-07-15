package common

import (
	"testing"
)

func Test_RandomString(t *testing.T) {
	input := []int{6, 10, 20, 16, 32}
	for i := range input {
		l := input[i]
		out := RandomString(l)
		if len(out) != l {
			t.Fatal("length not equal want :", l, " out:", out)
		}
	}
}

func Benchmark_RandomString(b *testing.B) {
	input := []int{10, 20, 16, 32}
	l := 0
	out := ""
	for i := 0; i < b.N; i++ {
		l = i % 4
		l = input[l]
		out = RandomString(l)
		if len(out) != l {
			b.Fatal("length not equal want :", l, " out:", out)
		}
	}
}
