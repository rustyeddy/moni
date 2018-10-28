package moni

import (
	"io"
	"time"

	log "github.com/sirupsen/logrus"
)

var (
	config *Configuration
)

type ConfigLogger struct {
	LevelString  string
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

type Command struct {
	Command string
	Args    []string
}

var (
	DefaultConfig = Configuration{
		Addrport: ":8888",
		Daemon:   false,
		Depth:    3,
		Profile:  false,
		Pubdir:   "docs",
		NoServe:  false,
		Storedir: "/srv/moni", // or "./.moni"
		Tmpldir:  "tmpl",
	}
)

// SetConfig sets and reconfigures the application
func SetConfig(cfg *Configuration) {
	// Set the global config
	config = cfg

	if config.Logerr == nil {
		config.Logerr = NewLogerr("config")
	}
}

func StoreConfig() {
	col := mdb.Collection("config")
	IfNilError(col, "getting config collection")
	panic("todo store the configuration")
}

func FetchConfig() {
	col := mdb.Collection("config")
	IfNilError(col, "FetchConfig")

	panic("todo fetch the configuration")
}
