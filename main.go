package main

import (
	"flag"
	"fmt"
	"net/url"
	"os"

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
}

var (
	config  Configuration
	err     error
	acl     map[string]bool
	pages   map[url.URL]*Page
	storage *store.FileStore

	urlChan  chan string
	pageChan chan *Page
)

func init() {
	flag.StringVar(&config.Addrport, "addr", "0.0.0.0:1212", "Address and port configuration")
	flag.StringVar(&config.ConfigFile, "config", "moni.json", "Moni config file")
	flag.StringVar(&config.LogFile, "logfile", "", "Crawl logfile")
	flag.StringVar(&config.LogFormat, "format", "", "format to print [json]")
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

	urls := flag.Args()
	if urls == nil {
		fmt.Println(flag.Usage)
	}

	doneChan := make(chan bool)
	urlChan = make(chan string)
	go Pager(pageChan, doneChan)
	go Scrubber(urlChan, pageChan, doneChan)

	for _, ustr := range urls {
		fmt.Printf("Walking %s send to scrubber\n", ustr)
		urlChan <- ustr
	}
	<-doneChan
	log.Infoln("The end, good bye ... ")
}

func setupLogging() {
	if config.LogFile != "" {
		f, err := os.OpenFile(config.LogFile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
		errFatal(err, "Failed to open "+config.LogFile)
		log.SetOutput(f)
	}
	if config.LogFormat == "json" {
		log.SetFormatter(&log.JSONFormatter{})
	}
}

func setupStorage() {
	if storage, err = store.UseFileStore("."); err != nil || storage == nil {
		errFatal(err, "failed to useFileStore ")
	}
}

// errPanic something went wrong, panic.
func errPanic(err error) {
	if err != nil {
		panic(err)
	}
}

// nilPanic does so when it's parameter is such.
func nilPanic(val interface{}) {
	if val == nil {
		fmt.Printf("val is nil")
	}
}

// errPanic something went wrong, panic.
func errFatal(err error, str string) {
	if err != nil {
		log.Fatalln(err, str)
	}
}

// nilPanic does so when it's parameter is such.
func nilFatal(val interface{}, str string) {
	if val == nil {
		log.Fatalln(str)
	}
}
