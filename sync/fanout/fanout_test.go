package fanout

import (
	"fmt"
	"testing"
)

func Test_Merge(t *testing.T){
	c := gen(2, 3, 4, 5, 6, 7, 8)
	out2 := sq(c)
	out1 := sq(c)
	for v := range Merge(len(c),out1, out2) {
		fmt.Println(v)
	}
}
func gen(nums ...int) <-chan interface{} {
	out := make(chan interface{})
	go func() {
		for _, n := range nums {
			out <- n
		}
		close(out)
	}()
	return out
}
func sq(in <-chan interface{}) <-chan interface{} {
	out := make(chan interface{})
	go func() {
		for n := range in {
			out <- (n.(int))*(n.(int))
		}
		close(out)
	}()
	return out
}