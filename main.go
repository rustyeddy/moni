package main

/*
 * Moni is my website monitoring tool
 */

import (
	"flag"
	"fmt"
	"log"
	"net/url"

	"github.com/gocolly/colly"
	"github.com/rustyeddy/store"
)

// Configuration manages all variables and parameters for a given run of moni.
type Configuration struct {
	ConfigFile string
	Changed    bool
	Verbose    bool
}

var (
	storage *store.FileStore
	config  Configuration
	baseURL string
	pages   map[string]*Page
	acl     map[string]bool
)

func init() {
	var err error

	flag.StringVar(&config.ConfigFile, "config", "moni.json", "Moni config file")
	storage, err = store.UseFileStore(".")
	if err != nil {
		panic(err)
	}
	pages = make(map[string]*Page)
	acl = make(map[string]bool)
}

func main() {
	flag.Parse()

	urls := flag.Args()
	if urls == nil || len(urls) == 0 {
		log.Fatal("Expected some sites, got none")
	}

	for _, baseURL = range urls {
		// Place the command line url in the acl allowed list
		if config.Verbose {
			fmt.Print("Add website ", baseURL, " to ACL")
		}
		acl[baseURL] = true
		Crawl(baseURL)
	}

	// Save the configuration if it has changed
	if config.Changed {
		fmt.Println("Configuration has changed, writing config file")
		storage.Save("config.json", config)
	}
	fmt.Println("The end...")
}

func errPanic(err error) {
	if err != nil {
		panic(err)
	}
}

func processPage(urlstr string) bool {
	var ok, ex bool

	u, err := url.Parse(urlstr)
	errPanic(err)

	if u.Hostname() == "" {
		// we will accept relative urls because the are belong to
		// the website being searched.
		return ok

	} else {
		if config.Verbose {
			fmt.Printf("Hostname %s\n", u.Hostname())
		}
	}

	if ok, ex = acl[urlstr]; ex == false {
		return false
	}

	if _, exists := pages[urlstr]; exists {
		return false
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
	if processPage(urlstr) {
		fmt.Println("\turl to be processed: ", urlstr)
		e.Request.Visit(urlstr)
	}
}

func doRequest(r *colly.Request) {
	fmt.Println("Request ", r.URL)
}
