package moni

import (
	"io"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

type ConfigLogger struct {
	LogLevel     string
	FormatString string
	Logfile      string

	log.Level
	log.Formatter
	Out io.Writer
}

type Configuration struct {
	Addrport string
	Basedir  string
	Client   bool

	ConfigLogger

	ConfigFile string
	Cli        bool
	Daemon     bool
	Debug      bool
	Depth      int

	NoServe bool

	Profile bool
	Pubdir  string

	Storedir string
	Tmpldir  string
	Wait     time.Duration

	*Logerr
}

// DefaultConfig sets some reasonable defaults before any
// configs can be read or any input from the external env
func DefaultConfig() Configuration {
	return Configuration{
		Addrport: ":8888",
		Daemon:   false,
		Depth:    3,
		Profile:  false,
		Pubdir:   "docs",
		NoServe:  false,
		Storedir: "/srv/moni",
		Tmpldir:  "tmpl",
	}
}

// SetConfig sets and reconfigures the application
func SetConfig(cfg *Configuration) {
	config = cfg
	if config.Logerr == nil {
		config.Logerr = NewLogerr("config")
	}
}

// SetLogger adjusts logger configs after a configuration change, such as
// immediately after parsing flags
func (c *Configuration) SetLogger() {
	if config.Logfile != "stdout" && config.Logfile != "stderr" {
		if wr, err := os.Open(config.Logfile); err != nil {
			log.Fatalln("Failed to open logfile ", config.Logfile)
		} else {
			log.SetOutput(wr)
		}
	}

	if l, err := log.ParseLevel(config.LogLevel); err != nil {
		log.Fatalln("failed to parse level ", config.Level)
	} else {
		log.SetLevel(l)
	}

	switch config.FormatString {
	case "text":
		log.SetFormatter(&log.TextFormatter{})
	case "json":
		log.SetFormatter(&log.JSONFormatter{})
	default:
		log.Fatalln("can not use formatter ", config.FormatString)
	}
}

func (c *Configuration) DebugLogger() {
	log.SetFormatter(&log.TextFormatter{
		DisableColors:    true,
		DisableTimestamp: true,
		DisableSorting:   true,
	})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

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

func ProdLogger() (l *log.Logger) {
	l = log.New()
	l.SetFormatter(&log.JSONFormatter{})
	l.SetOutput(os.Stdout)
	l.SetLevel(log.WarnLevel)
	return l
}
