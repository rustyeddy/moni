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

type Logmap map[string]*Logerr

var (
	log    *Logerr
	logmap *Logmap
)

func init() {
	logmap = new(Logmap)
	log = newDevJSON("main")
}

func NewLogerr(name string) (nl *Logerr) {
	nl = &Logerr{Name: name}
	nl.Logger = logrus.New()
	return nl
}

// NewDebugger
func newDevtest(name string) (l *Logerr) {
	l = NewLogerr(name)
	l.SetFormatter(&logrus.TextFormatter{})
	l.SetOutput(os.Stdout)
	l.SetLevel(logrus.DebugLevel)
	return l
}

// NewDebugger
func newDevJSON(name string) (l *Logerr) {
	l = NewLogerr(name)
	l.SetFormatter(&logrus.JSONFormatter{})
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

// Clone an existing logger, with a new name
func (l *Logerr) Clone(name string) (nl *Logerr) {
	nl = NewLogerr(name)
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
