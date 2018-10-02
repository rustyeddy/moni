package inventory

import (
	"flag"
)

type Configuration struct {
	Output string
	Level  string
	Format string

	Addrport string
	Pubdir   string
	Depth    int
	Daemon   bool
}

var (
	Config Configuration
)

func init() {
	flag.StringVar(&Config.Output, "output", "stdout", "Were to send log output")
	flag.StringVar(&Config.Level, "level", "warn", "Log level to set")
	flag.StringVar(&Config.Format, "format", "json", "Format to print log files")

	flag.StringVar(&Config.Addrport, "addr", ":4444", "Run an Daemon in the background")
	flag.IntVar(&Config.Depth, "depth", 1, "Max crawl depth")
	flag.BoolVar(&Config.Daemon, "daemon", true, "Run an Daemon in the background")
	flag.StringVar(&Config.Pubdir, "dir", "./pub", "Run an Daemon in the background")
}

func GetConfiguration() *Configuration {
	return &Config // make a copy
}
