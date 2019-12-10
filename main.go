package main

import (
	"flag"
	"fmt"
	"net/url"

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
	storage *store.FileStore

	urlChan chan string
)

func init() {
	flag.StringVar(&config.Addrport, "addr", "0.0.0.0:1212", "Address and port configuration")
	flag.StringVar(&config.ConfigFile, "config", "crawl.json", "Moni config file")
	flag.StringVar(&config.LogFile, "logfile", "crawl.log", "Crawl logfile")
	flag.StringVar(&config.LogFormat, "format", "", "format to print [json]")
	flag.StringVar(&config.Pubdir, "pub", "pub", "the published dir")
	flag.BoolVar(&config.Recurse, "recurse", true, "Recurse local")

	//storage, err := store.UseFileStore(".")
	//errPanic(err)
	pages = make(map[url.URL]*Page)
	acl = make(map[string]bool)
	urlChan = make(chan string)

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

	doneChan := make(chan bool)
	urlChan = make(chan string)
	go Scrubber(urlChan, doneChan)
	go doRouter(config.Pubdir)

	// Process whatever was passed to us
	processURLs(flag.Args())

	<-doneChan
	log.Infoln("The end, good bye ... ")
}

func processURLs(urls []string) {
	for _, ustr := range urls {
		fmt.Printf("Walking %s\n", ustr)
		urlChan <- ustr
	}
}

func setupStorage() {
	if storage, err = store.UseFileStore("."); err != nil || storage == nil {
		errFatal(err, "failed to useFileStore ")
	}
}
