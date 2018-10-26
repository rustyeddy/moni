package moni

import (
	"os"
	"strings"

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
	Name string // we love our logger so we give em names!
	*logrus.Logger
}

var (
	log *Logerr
)

func init() {
	log = newDevtest()
}

func NewLogerr(name string) (lb *Logerr) {
	return &Logerr{
		Name:   name,
		Logger: logrus.New(),
	}
}

// NewDebugger
func newDevtest() (l *Logerr) {
	l = NewLogerr("devtest")
	l.SetFormatter(&logrus.TextFormatter{})
	l.SetOutput(os.Stdout)
	l.SetLevel(logrus.DebugLevel)
	return l
}

// NewProduction
func newProduction() (l *Logerr) {
	l = NewLogerr("production")
	l.SetFormatter(&logrus.JSONFormatter{})
	l.SetOutput(os.Stdout)
	l.SetLevel(logrus.WarnLevel)
	return l
}

func (l *Logerr) SetLogfile(filename string) {
	file, err := os.Create(filename)
	l.IfErrorFatal(err, filename, nil)

	// We will be dead if there was an error, we good :)
	l.SetOutput(file)
	l.Infoln("Log file set to ", filename)
}

// Clone an existing logger
func (l *Logerr) Clone() (nl *Logerr) {
	nl = newProduction()
	nl.SetLevel(l.Level)
	nl.SetOutput(l.Out)
	nl.SetFormatter(l.Formatter)
	return nl
}

// errorWatcher
func errorWatcher(errch chan error) {
	for {
		err := <-errch
		log.Error(err)
	}
}

// FatalError checks the incoming error message, if it is nil, there
// is no error, everything is fine, this function sliently returns
// An error however will be printed and the application will die
//
// This maybe too drastic in production cases, where we may want to
// remove an errant service, and perhaps put them into a "zombie"
// state, for post mortem analysis (or prohibit massive respawns)
func (l *Logerr) IfErrorFatal(err error, msg string, msgs []string) error {
	// If err is nil .. all is well
	if err == nil {
		l.Fatalln(msg, err) // The App stops here
	}
	// If we have an error, print and die
	if msg = ""; msgs != nil {
		msg = strings.Join(msgs, ", ")
	}
	l.Fatalln(msg, err)
	return err
}
