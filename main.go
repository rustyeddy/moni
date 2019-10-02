package main

/*
 * Moni is my website monitoring tool
 */

import (
	"flag"
	"fmt"
	"log"

	"github.com/rustyeddy/store"
)

// Configuration manages all variables and parameters for a given run of moni.
type Configuration struct {
	ConfigFile string
}

var (
	storage *store.FileStore
	config  Configuration
)

func init() {
	var err error

	flag.StringVar(&config.ConfigFile, "config", "moni.json", "Moni config file")
	storage, err = store.UseFileStore(".")
	if err != nil {
		panic(err)
	}
}

func main() {
	flag.Parse()

	urls := flag.Args()
	if urls == nil || len(urls) == 0 {
		log.Fatal("Expected some sites, got none")
	}

	for _, u := range urls {
		site := NewSite(u)
		site.Crawl()
	}
	fmt.Println("The end...")
}

func err_panic(err error) {
	if err != nil {
		panic(err)
	}
}
