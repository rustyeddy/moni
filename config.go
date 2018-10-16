package main

import (
	"errors"
	"flag"
	"time"

	log "github.com/sirupsen/logrus"
)

type ConfigLogger struct {
	Output string
	Level  string
	Format string
}

type Configuration struct {
	Addrport string
	Client   bool

	ConfigLogger
	ConfigFile string
	Cli        bool
	Daemon     bool

	Depth int

	Profile bool
	Pubdir  string // Where to serve the static files from

	Serve    bool
	StoreDir string
	Wait     time.Duration
}

type Command struct {
	Command string
	Args    []string
}

func init() {
	flag.StringVar(&Config.Output, "output", "stdout", "Were to send log output")
	flag.StringVar(&Config.Level, "level", "warn", "Log level to set")
	flag.StringVar(&Config.Format, "format", "json", "Format to print log files")

	flag.StringVar(&Config.Addrport, "addr", ":8888", " an Daemon in the background")

	// What "cmd" or "mode" to run the command crawl, run cli or daemon
	flag.BoolVar(&Config.Cli, "cli", false, "Run a command line client")
	flag.BoolVar(&Config.Serve, "serve", false, "Run as a service")

	flag.StringVar(&Config.ConfigFile, "cfg", "/srv/inv/config.json", "Use configuration file")

	flag.IntVar(&Config.Depth, "depth", 1, "Max crawl depth")
	flag.StringVar(&Config.Pubdir, "dir", "pub", "Serve the site from this dir")
	flag.StringVar(&Config.StoreDir, "store", "/srv/inv", "Directory for Store to use")

	flag.BoolVar(&Config.Profile, "prof", false, "Profile our http server (daemon)")
}

func GetConfiguration() *Configuration {
	return &Config
}

// SaveFile will write the configuration out to our Storage as JSON.  As of
// now we are using the Storage interface to save our configurations.  Hence,
// we will use a label name and not a filename, although the label name can
// look like a filename..  That is a label can be "config.json" but it can
// NOT have any leading path '/' characters
func (c *Configuration) SaveFile() {
	s := getStorage()
	if s == nil {
		log.Errorln("Failed saving config ", c.ConfigFile)
		return
	}
	_, err := s.StoreObject("config", c)
	if err != nil {
		log.Errorln("Failed writing configuration", c.ConfigFile, err)
	}
}

// ReadFile fetches our configuration object from our storage container,
// if needed, the object will be converted from JSON to a Go object before
// being returned.
func (c *Configuration) ReadFile() error {
	s := getStorage() // fatal if nul
	_, err := s.FetchObject("config", c)
	if err != nil {
		return errors.New("failed to read config from store")
	}
	return nil
}
