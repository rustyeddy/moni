package main

type MoniError struct {
	msg string
}

var (
	ErrorNil      MoniError = MoniError{"expected (obj) got ()"}
	ErrorNotFound MoniError = MoniError{"not found"}
)

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
