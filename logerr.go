package moni

import (
	"io"
	"os"
	"strings"

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

type Logerr struct {
	Name string // we love our logger so we give em names!
	*log.Logger
}

func NewLogerr(name string) (nl *Logerr) {
	return &Logerr{Name: name}
}

func SetValues(out io.Writer, formatter log.Formatter, level log.Level) {
	log.SetOutput(out)
	log.SetFormatter(formatter)
	log.SetLevel(level)
}

// SetValues is simple to set ofl values
func (l *Logerr) SetValues(out io.Writer, formatter log.Formatter, level log.Level) *Logerr {
	l.SetFormatter(formatter)
	l.SetOutput(out)
	l.SetLevel(level)
	return l
}

func (l *Logerr) SetDebugging() {
	l.SetValues(os.Stdout, &log.TextFormatter{}, log.DebugLevel)
}

func (l *Logerr) SetTesting() {
	l.SetValues(os.Stdout, &log.JSONFormatter{}, log.WarnLevel)
}

func (l *Logerr) SetProduction(filename string) {
	file, err := os.Create(filename)
	IfErrorFatal(err, "SetTesting", filename)
	l.SetValues(file, &log.JSONFormatter{}, log.WarnLevel)
}

// errorWatcher = log error messages
func (lm *Logerr) WatchChannel(errch chan error) {
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
func IfErrorFatal(err error, msgs ...string) error {
	// If err is nil .. all is well
	if err == nil {
		return nil // we are good, nothing to do
	}
	// If we have an error, print and die
	msg := ""
	if msgs != nil {
		msg = strings.Join(msgs, ", ")
	}
	log.Fatalln(msg, err)
	return err
}

// FatalError checks the incoming error message, if it is nil, there
// is no error, everything is fine, this function sliently returns
// An error however will be printed and the application will die
//
// This maybe too drastic in production cases, where we may want to
// remove an errant service, and perhaps put them into a "zombie"
// state, for post mortem analysis (or prohibit massive respawns)
func IfNilError(obj interface{}, msgs ...string) bool {
	// If err is nil .. all is well
	if obj != nil {
		return false // we are good, nothing to do
	}

	// If we have an error, print and die
	msg := ""
	if msgs != nil {
		msg = strings.Join(msgs, ", ")
	}
	log.Errorln(msg)
	return true
}

// ====================================================================

// Logmap is the struct holding loggers in the event our modules want
// to share the loggers (turn on / off and change params all at once)
type Logmap map[string]*Logerr

// AddLogger will add the given logger with name to the logmap. If the
// logger already exists, it will be overwritten.
func (lm Logmap) AddLogger(name string, lg *Logerr) Logmap {
	lm[name] = lg
	return lm
}

type LogManager struct {
	Logmap
	LogChannel chan Err
}

// Errors
// ====================================================================
type Err struct {
	errno int
	error
}
