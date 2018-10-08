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
	ConfigLogger
	Pubdir     string // Where to serve the static files from
	Depth      int
	Addrport   string
	StoreDir   string
	Client     bool
	Profile    bool
	ConfigFile string
	Wait       time.Duration
}

func init() {
	flag.StringVar(&Config.Output, "output", "stdout", "Were to send log output")
	flag.StringVar(&Config.Level, "level", "warn", "Log level to set")
	flag.StringVar(&Config.Format, "format", "json", "Format to print log files")
	flag.StringVar(&Config.Addrport, "http-addr", ":8888", " an Daemon in the background")
	flag.IntVar(&Config.Depth, "depth", 1, "Max crawl depth")
	flag.BoolVar(&Config.Client, "cli", false, "Run a command line client")
	flag.StringVar(&Config.Pubdir, "dir", "./pub", "Run an Daemon in the background")
	flag.StringVar(&Config.StoreDir, "store", "/srv/inv", "Directory for Store to use")
	flag.StringVar(&Config.ConfigFile, "cfg", "/srv/inv/config.json", "Use configuration file")
	flag.DurationVar(&Config.Wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")

	flag.BoolVar(&Config.Profile, "prof", false, "Profile our http server")
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
