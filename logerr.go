package moni

import (
	"os"

	"github.com/sirupsen/logrus"
)

const (
	ErrNone = 0
	ErrGeneric
	ErrNotFound
	ErrNotSupported
	ErrInternalError
	ErrNilObject
)

type Logerr struct {
	*logrus.Logger
}

var (
	log *Logerr
)

func init() {
	log = newDevtest()
}

// NewDebugger
func newDevtest() (l *Logerr) {
	l = &Logerr{logrus.New()}
	l.SetFormatter(&logrus.TextFormatter{})
	l.SetOutput(os.Stdout)
	l.SetLevel(logrus.DebugLevel)
	return l
}

// NewProduction
func newProduction() (l *Logerr) {
	l = &Logerr{logrus.New()}
	l.SetFormatter(&logrus.JSONFormatter{})
	l.SetOutput(os.Stdout)
	l.SetLevel(logrus.WarnLevel)
	return l
}

// errorWatcher
func errorWatcher(errch chan error) {
	for {
		err := <-errch
		log.Error(err)
	}
}
