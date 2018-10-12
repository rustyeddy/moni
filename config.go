package main

import (
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

func (c *Configuration) SaveFile() {
	_, err := Storage.StoreObject(c.ConfigFile, c)
	if err != nil {
		log.Errorln("Failed writing configuration", c.ConfigFile, err)
	}
}
