package main

import "errors"

type MoniError struct {
	msg string
}

var (
	ErrorNil      MoniError = MoniError{"expected (obj) got ()"}
	ErrorNotFound MoniError = MoniError{"not found"}
)

func ErrorNotSupported(msg string, err error) error {
	if err != nil {
		msg += " " + err.Error()
	}
	return errors.New(msg)
}

// Error returns the error message and satisfies the error.Error interface
func (m MoniError) Error() string {
	return m.msg
}

func AssertNotNil(obj interface{}) error {
	if obj == nil {
		err := ErrorNotFound
		return err
	}
	return nil
}
