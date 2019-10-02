package moni

import (
	"io"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

// ConfigLogger controlls how the logger is configured.
type ConfigLogger struct {
	LogLevel     string
	FormatString string
	Logfile      string

	log.Level
	log.Formatter
	Out io.Writer
}

// Configuration sports the major configurable items
// of our moni application
type Configuration struct {
	Addrport string
	Basedir  string

	ConfigLogger
	ConfigFile string

	Profile bool
	Debug   bool
	Depth   int

	Pubdir   string
	Tmpldir  string
	Storedir string

	Wait time.Duration
	*Logerr
}

// DefaultConfig sets some reasonable defaults before any
// configs can be read or any input from the external env
func DefaultConfig() Configuration {
	return Configuration{
		Addrport: ":8888",
		Depth:    3,
		Pubdir:   "docs",
		Storedir: "/srv/moni",
		Tmpldir:  "tmpl",
	}
}

// SetLogger adjusts logger configs after a Configuration change, such as
// immediately after parsing flags
func (c *Configuration) SetLogger() {
	if app.Logfile != "stdout" && app.Logfile != "stderr" {
		if wr, err := os.Open(app.Logfile); err != nil {
			log.Fatalln("Failed to open logfile ", app.Logfile)
		} else {
			log.SetOutput(wr)
		}
	}

	if l, err := log.ParseLevel(app.LogLevel); err != nil {
		log.Fatalln("failed to parse level ", app.Level)
	} else {
		log.SetLevel(l)
	}

	switch app.FormatString {
	case "text":
		log.SetFormatter(&log.TextFormatter{})
	case "json":
		log.SetFormatter(&log.JSONFormatter{})
	default:
		log.Fatalln("can not use formatter ", app.FormatString)
	}
}

// DebugLogger returns a logger ready for debugging
func (c *Configuration) DebugLogger() {
	log.SetFormatter(&log.TextFormatter{
		DisableColors:    true,
		DisableTimestamp: true,
		DisableSorting:   true,
	})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

// Write the contents of the configuration file.
func (c *Configuration) Write(bufer []byte) (err error) {
	panic("TODO")
	return err
}

// GetDebugLogger returns a logger ready for debugging
func GetDebugLogger() (l *log.Logger) {
	l = log.New()
	l.SetFormatter(&log.TextFormatter{
		DisableColors:    false,
		DisableTimestamp: true,
		DisableSorting:   true,
	})
	l.SetOutput(os.Stdout)
	l.SetLevel(log.DebugLevel)
	return l
}

// ProdLogger is a production logger
func ProdLogger() (l *log.Logger) {
	l = log.New()
	l.SetFormatter(&log.JSONFormatter{})
	l.SetOutput(os.Stdout)
	l.SetLevel(log.WarnLevel)
	return l
}
