package common

import (
	"errors"
	"fmt"
	"testing"
)

func Test_Error(t *testing.T) {
	e := NewError(1, fmt.Errorf("%v", "some error"))
	t.Log(e, e.Code)
	t.Logf("%s", e)

	emsg := NewErrorWithMsg(2, "some error")
	t.Log(emsg, emsg.Code)
}

func Test_AssertError(t *testing.T) {
	var targetErr = NewError(1, fmt.Errorf("%v", "some error"))
	var e error = targetErr
	if !errors.Is(e, targetErr) {
		t.Fatal("errors.Is not equal")
	}
	if errors.Unwrap(e) == nil {
		t.Fatal("errors.Unwrap not nil")
	}
	var commErr Error
	if !errors.As(e, &commErr) {
		t.Fatal("errors.As error")
	}
}
