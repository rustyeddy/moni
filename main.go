package main

import (
	"flag"
	"sync"

	"github.com/rustyeddy/store"
	log "github.com/sirupsen/logrus"
)

// Configuration manages all variables and parameters for a given run of moni.
type Configuration struct {
	Addrport   string
	ConfigFile string
	Changed    bool
	Daemon     bool
	LogFile    string
	LogFormat  string
	Pubdir     string
	Recurse    bool
	Verbose    bool
	Wait       int
}

var (
	config  Configuration
	err     error
	acl     map[string]bool
	sites   Sites
	storage *store.FileStore

	walkQ chan *Page
)

func init() {
	flag.StringVar(&config.Addrport, "addr", "0.0.0.0:2222", "Address and port configuration")
	flag.StringVar(&config.ConfigFile, "config", "crawl.json", "Moni config file")
	flag.StringVar(&config.LogFile, "logfile", "", "Crawl logfile")
	flag.StringVar(&config.LogFormat, "format", "", "format to print [json]")
	flag.StringVar(&config.Pubdir, "pub", "pub", "the published dir")
	flag.BoolVar(&config.Recurse, "recurse", true, "Recurse local")
	flag.BoolVar(&config.Daemon, "daemon", true, "Run as a service opening and listening to sockets")
	flag.BoolVar(&config.Verbose, "verbose", false, "turn on or off verbosity")
	flag.IntVar(&config.Wait, "wait", 5, "wait in minutes between check")

	//storage, err := store.UseFileStore(".")
	//errPanic(err)
	acl = make(map[string]bool)

	sites = make(Sites)
	walkQ = make(chan *Page, 100)

	// TODO read the acls from a file
	acl["localhost"] = false
	acl["google.com"] = false
	acl["github.com"] = false
	acl["rustyeddy.com"] = true
}

func main() {
	var wg sync.WaitGroup

	flag.Parse()
	setupLogging()
	setupStorage()

	wg.Add(2)
	go doRouter(config.Pubdir, &wg)
	go doWalker(walkQ, &wg)

	setupSites(walkQ, flag.Args())
	wg.Wait()

	log.Infoln("The end, good bye ... ")
}
