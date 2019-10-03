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
	Changed bool
}

var (
	storage *store.FileStore
	config  Configuration
	pages	map[string]*Page
	acl		map[string]bool
)

func init() {
	var err error

	flag.StringVar(&config.ConfigFile, "config", "moni.json", "Moni config file")
	storage, err = store.UseFileStore(".")
	if err != nil {
		panic(err)
	}
	pages = make(map[string]*Page)
}

func main() {
	flag.Parse()

	urls := flag.Args()
	if urls == nil || len(urls) == 0 {
		log.Fatal("Expected some sites, got none")
	}

	for _, u := range urls {
		Crawl(u)
	}
	fmt.Println("The end...")

	// Save the configuration if it has changed
	if config.Changed {
		fmt.Println("Configuration has changed, writing config file")
		storage.Save("config.json", config)		
	}
}

func panic_err(err error) {
	if err != nil {
		panic(err)
	}
}

func processPage(urlstr string) bool {
	var ok, ex bool
	
	if _, exists := pages[urlstr]; exists == false {
		if ok, ex = acl[urlstr]; ex == false {
			return false
		}
	}
	return ok
}

// Crawl the given URL
func Crawl(urlstr string) {
	c := colly.NewCollector()

	// Setup all the collbacks
	c.OnHTML("a", doHyperlink)
	c.OnRequest(doRequest)
	c.Visit(urlstr)
}

func doHyperlink(e *colly.HTMLElement) {
	urlstr := e.Attr("href")
	fmt.Println("url: ", urlstr)
	//if processPage(urlstr) {
	fmt.Println("\turl to be processed: ", urlstr)
	e.Request.Visit(urlstr)
	//}
}

func doRequest(r *colly.Request) {
	fmt.Println("Request ", r.URL)
}
