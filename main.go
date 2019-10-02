package main

/*
 * Moni is my website monitoring tool
 */

import (
	"flag"
	"fmt"
	"log"

	"github.com/gocolly/colly"
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

	sites := flag.Args()
	if sites == nil || len(sites) == 0 {
		log.Fatal("Expected some sites, got none")
	}

	for _, s := range sites {
		Crawl(s)
	}
	fmt.Println("The end...")
}

// Crawl the given URL
func Crawl(baseurl string) {
	c := colly.NewCollector()

	// Find and visit all links
	c.OnHTML("a", func(e *colly.HTMLElement) {
		thing := e.Attr("href")

		page, err := processPage(thing)
		err_panic(err)
		
		if page != nil {
			e.Request.Visit(thing)
		}
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	fmt.Println("visiting", baseurl)
	c.Visit(baseurl)
	storage.Save("config.json", config)
}

func err_panic(err error) {
	if err != nil {
		panic(err)
	}
}
