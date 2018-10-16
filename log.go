package main

import "fmt"

const (
	NilObjectError = iota
	NotSupportedError
	NotFoundError
)

type Error struct {
	errnum int
	msg    string
}

var (
	errorNil          = &Error{NilObjectError, "expected (obj) got ()"}
	errorNotSupported = &Error{NotSupportedError, "not supported"}
)

// Error returns the error message and satisfies the error.Error interface
func (e Error) Error() string {
	return e.msg
}

func (e Error) String() string {
	return fmt.Sprintf("Error(%d) %s", e.errnum, e.msg)
}

func ErrorNotSupported(msg string) *Error {
	e := errorNotSupported
	e.msg = msg
	return e
}

func ErrorNil(msg string) *Error {
	e := errorNil
	e.msg = msg
	return e
}

func AssertNotNil(obj interface{}) *Error {
	if obj == nil {
		err := errorNil
		return err
	}
	return nil
}
