package inv

import (
	"flag"
)

type ConfigLogger struct {
	Output string
	Level  string
	Format string
}

type Configuration struct {
	Log      ConfigLogger
	Addrport string
	Pubdir   string
	Depth    int
	Client   bool
}

var (
	Config Configuration
)

func init() {
	flag.StringVar(&Config.Log.Output, "output", "stdout", "Were to send log output")
	flag.StringVar(&Config.Log.Level, "level", "warn", "Log level to set")
	flag.StringVar(&Config.Log.Format, "format", "json", "Format to print log files")

	flag.StringVar(&Config.Addrport, "addr", ":4444", "Run an Daemon in the background")
	flag.IntVar(&Config.Depth, "depth", 1, "Max crawl depth")
	flag.BoolVar(&Config.Client, "cli", false, "Run a command line client")
	flag.StringVar(&Config.Pubdir, "dir", "./pub", "Run an Daemon in the background")
}

func GetConfiguration() *Configuration {
	return &Config
}
