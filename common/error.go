package common

import "fmt"

type Error struct {
	Code int
	err  error
}

func NewError(code int, e error) Error {
	return Error{
		Code: code,
		err:  e,
	}
}

func NewErrorWithMsg(code int, msg string) Error {
	return NewError(code, fmt.Errorf(msg))
}

func (e Error) Error() string {
	if e.err == nil {
		return ""
	}
	return e.err.Error()
}

func (e Error) Unwrap() error {
	return e.err
}
