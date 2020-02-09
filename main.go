package main

import (
	"flag"

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
	Wait       int
}

// Global variables are all declared here, keeps them in
// one easy to track spot.
var (
	config  Configuration    // Configuration from above
	acl     ACL              // ACL controls which URLs are crawled or blocked
	sites   Sites            // The map of Sites we are watching
	storage *store.FileStore // Local file storage TODO: add DO and GCP
	walkQ   chan *Page       // channel used to submit pages (or urls) to be walked
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
	flag.IntVar(&config.Wait, "wait", 5, "wait in minutes between check")

	sites = make(Sites)
	walkQ = make(chan *Page, 100)
}

func main() {
	flag.Parse()
	setupLogging()
	setupStorage()

	slist := readSitesFile()
	slist = append(flag.Args())

	// scrubSites returns a channel from a go routine producing *Page
	// objects from raw, unscrubed URLs. If the URL is rejected it will
	// not be returned as a page structure.
	for page := range scrubSites(slist) {
		page.Walk()
	}
	log.Infoln("The end, good bye ... ")
}
