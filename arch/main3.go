package moni

import (
	"flag"
)

var (
	cfg moni.Configuration
)

func init() {

	cfg = moni.DefaultConfig()

	// Logerr settings
	flag.StringVar(&cfg.Logfile, "logfile", "stdout", "Were to send log output")
	flag.StringVar(&cfg.LogLevel, "level", "warn", "Log level to set")
	flag.StringVar(&cfg.FormatString, "format", "json", "Format to print log files")

	// Crank up the verbosity a bit
	flag.BoolVar(&cfg.Debug, "debug", false, "Debug stuff with this program")

	// Address and Port the server to open and listen for connections
	flag.StringVar(&cfg.Addrport, "addr", ":8888", " an Daemon in the background")

	flag.StringVar(&cfg.ConfigFile, "cfg", "/srv/moni/cfg.json", "Use cfguration file")

	// -crawl related
	flag.IntVar(&cfg.Depth, "crawl-depth", 1, "Max crawl depth")
	flag.StringVar(&cfg.Pubdir, "cfgdir", "pub", "Serve the site from this dir")
	flag.StringVar(&cfg.Storedir, "storedir", "/srv/moni/", "Directory for Store to use")
	flag.StringVar(&cfg.Tmpldir, "tmpldir", "../tmpl", "Basedir for templates")

	flag.BoolVar(&cfg.Profile, "prof", true, "Profile our http server (daemon)")
}

func main3() {
	// Flags are mostly set in the moni.config.go package
	flag.Parse()

	// Gets the cfg, and saves the configuration file with the cfg
	app := moni.GetApp(&cfg)
	app.Init()
	app.StartService()
}
