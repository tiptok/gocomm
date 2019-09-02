package counter

import (
	"fmt"
	"testing"
)
type cc int64
func TestCounter(t *testing.T){
	var c cc
	c = 1
	fmt.Println(c)
	d :=2
	fmt.Println(d)
	d = int(c)
	fmt.Println(d)
}
