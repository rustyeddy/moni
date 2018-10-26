package moni

import (
	log "github.com/sirupsen/logrus"
)

const (
	ErrNone = 0
	ErrGeneric
	ErrNotFound
	ErrNotSupported
	ErrInternalError
	ErrNilObject
)

// errorWatcher
func errorWatcher(errch chan error) {
	for {
		err := <-errch
		log.Error(err)
	}
}
