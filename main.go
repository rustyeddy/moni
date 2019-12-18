package main

import (
	"flag"
	"net/url"
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
	Recurse    bool
	Verbose    bool
	LogFile    string
	LogFormat  string
	Pubdir     string
}

var (
	config  Configuration
	err     error
	acl     map[string]bool
	pages   map[url.URL]*Page
	sites   []string
	storage *store.FileStore

	doneChan chan bool
)

func init() {
	flag.StringVar(&config.Addrport, "addr", "0.0.0.0:1212", "Address and port configuration")
	flag.StringVar(&config.ConfigFile, "config", "crawl.json", "Moni config file")
	flag.StringVar(&config.LogFile, "logfile", "crawl.log", "Crawl logfile")
	flag.StringVar(&config.LogFormat, "format", "", "format to print [json]")
	flag.StringVar(&config.Pubdir, "pub", "pub", "the published dir")
	flag.BoolVar(&config.Recurse, "recurse", true, "Recurse local")
	flag.BoolVar(&config.Daemon, "daemon", false, "format to print [json]")
	flag.BoolVar(&config.Verbose, "verbose", false, "turn on or off verbosity")

	//storage, err := store.UseFileStore(".")
	//errPanic(err)
	acl = make(map[string]bool)
	doneChan = make(chan bool)
	pages = make(map[url.URL]*Page)
	sites = nil

	// TODO read the acls from a file
	acl["localhost"] = false
	acl["google.com"] = false
	acl["github.com"] = false
	acl["rustyeddy.com"] = true
}

func main() {
	flag.Parse()
	setupLogging()
	setupStorage()

	// Process urls will filter bad, redundant and blocked URLs
	// the urls that do not get blocked are then walked.
	if urls := flag.Args(); urls != nil && len(urls) > 0 {
		processURLs(urls, nil)
	}

	if config.Daemon {
		var wg sync.WaitGroup
		wg.Add(2)

		// Start the scrubber, router
		go doRouter(config.Pubdir, &wg)

		sites := GetSites()
		go watchSites(sites, &wg)

		wg.Wait()
	}

	log.Infoln("The end, good bye ... ")
}
