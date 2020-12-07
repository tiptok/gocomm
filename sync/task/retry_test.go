package task

import (
	"errors"
	"testing"
)

func TestDoWithRetries(t *testing.T) {
	count := 0
	DoWithRetries(func() error {
		count++
		return errors.New("retry error")
	}, WithRetries(5))
	if count != 5 {
		t.Fatal("except retry time:", 5, "actual time:", count)
	}
	//t.Log(err.Error())
}
