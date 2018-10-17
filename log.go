package main

import (
	"fmt"
	"strings"
	"time"
)

/*

TODO:

- Create standard log messages and error strings
- Use fields provided by logrus and create some category specific loggers
- Systematize log messages to make analyzing them easier
*/

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

// GetTimeStamp returns a timestamp in a modified RFC3339
// format, basically remove all colons ':' from filename, since
// they have a specific use with Unix pathnames, hence must be
// escaped when used in a filename.
func TimeStamp() string {
	ts := time.Now().UTC().Format(time.RFC3339)
	return strings.Replace(ts, ":", "", -1) // get rid of offesnive colons
}
