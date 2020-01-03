package main

import (
	"flag"
	"sync"

	"github.com/rustyeddy/store"
	log "github.com/sirupsen/logrus"
)

// Configuration manages all variables and parameters that
// can be set at the beginning of the program, or changed
// during the programs run time.  These variables will change
// the behavior and outcome of the program.
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
	Wait       int64
}

// Global variables are all declared here, keeps them in
// one easy to track spot.
var (
	config  Configuration    // Configuration from above
	acl     ACL              // ACL controls which URLs are crawled or blocked
	sites   Sites            // The map of Sites we are watching
	storage *store.FileStore // Local file storage TODO: add DO and GCP
	walkQ   chan string      // the string needs to be parsed to URL and OKd
)

// Initialize the config file with defaults, and setup the program to accept
// the respective command line flags and arguments.  We also initialize the
// sites variable with an empty map.  We establish the walkQ channel with a
// backlog of X, we'll see if we need to adjust.
func init() {
	flag.StringVar(&config.Addrport, "addr", "0.0.0.0:2222", "Address and port configuration")
	flag.StringVar(&config.ConfigFile, "config", "crawl.json", "Moni config file")
	flag.StringVar(&config.LogFile, "logfile", "", "Crawl logfile")
	flag.StringVar(&config.LogFormat, "format", "", "format to print [json]")
	flag.StringVar(&config.Pubdir, "pub", "pub", "the published dir")
	flag.BoolVar(&config.Recurse, "recurse", true, "Recurse local")
	flag.BoolVar(&config.Daemon, "daemon", true, "Run as a service opening and listening to sockets")
	flag.BoolVar(&config.Verbose, "verbose", false, "turn on or off verbosity")
	flag.Int64Var(&config.Wait, "wait", 5, "wait in minutes between check")

	sites = make(Sites)
	walkQ = make(chan string, 100)
}

func main() {
	var wg sync.WaitGroup

	// Parse command line arguments
	flag.Parse()
	setupLogging()
	setupStorage()

	wg.Add(2)
	go startRouter(config.Pubdir, &wg)
	go startWatcher(walkQ, &wg)

	slist := readSitesFile()
	slist = append(flag.Args())

	submitSites(slist)
	wg.Wait()

	log.Infoln("The end, good bye ... ")
}
