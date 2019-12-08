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
	LogFile    string
	Recurse    bool
	Verbose    bool
}

var (
	config  Configuration
	err     error
	acl     map[string]bool
	pages   map[url.URL]*Page
	storage *store.FileStore
)

func init() {
	flag.StringVar(&config.Addrport, "addr", "0.0.0.0:1212", "Address and port configuration")
	flag.StringVar(&config.ConfigFile, "config", "moni.json", "Moni config file")
	flag.StringVar(&config.LogFile, "logfile", "crawl-log.json", "Crawl logfile")
	flag.BoolVar(&config.Recurse, "recurse", true, "Recurse local")

	//storage, err := store.UseFileStore(".")
	//errPanic(err)
	pages = make(map[url.URL]*Page)
	acl = make(map[string]bool)

	// TODO read the acls from a file
	acl["localhost"] = false

	log.SetFormatter(&log.JSONFormatter{})

}

func main() {
	flag.Parse()

	setupLogging()
	setupStorage()

	// Process URLs on the command line if we have them
	if flag.Args() != nil {
		scrubURLs(flag.Args())
	}
}

func setupLogging() {
	if config.LogFile != "" {
		f, err := os.OpenFile(config.LogFile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
		errFatal(err, "Failed to open "+config.LogFile)
		log.SetOutput(f)
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
