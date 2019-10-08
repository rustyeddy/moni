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

	// walk the command line arguments treating them as URLs
	for _, baseURL = range urls {
		// Place the command line url in the acl allowed list
		if config.Verbose {
			fmt.Print("Add website ", baseURL, " to ACL")
		}

		u, err := url.Parse(baseURL)
		errPanic(err)

		acl[u.Hostname()] = true

		// This will become sending a message
		Crawl(u)
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

	host := u.Hostname()
	if host == "" {
		// we will accept relative urls because the are belong to
		// the website being searched.
		return true
	} else {
		fmt.Printf("Hostname %s\n", host)
	}

	if ok, ex = acl[host]; ex {
		return ok
	}

	return false
}

// Crawl the given URL
func Crawl(u *url.URL) {
	c := colly.NewCollector()

	// Setup all the collbacks
	c.OnHTML("a", doHTML)
	c.OnRequest(doRequest)
	c.OnResponse(doResponse)
	c.OnScraped(doScraped)

	c.Visit(u.String())
}

func doHTML(e *colly.HTMLElement) {
	urlstr := e.Attr("href")
	u, _ := url.Parse(urlstr)

	fmt.Printf("url: %+v\n", u)
	if processPage(urlstr) {
		fmt.Println(" ok ...")
		e.Request.Visit(urlstr)
	} else {
		fmt.Println(" blocked ...")
	}
}

// Called before the request is sent
func doRequest(r *colly.Request) {
	fmt.Println("Request ", r.URL)
}

// Called after the response is
func doResponse(r *colly.Response) {
	fmt.Println("Response ", r.Request.URL)
}

func doScraped(r *colly.Response) {
	fmt.Println("Scraped ", r.Request.URL)
}

func doError(_ *colly.Response, err error) {
	fmt.Println("Error", err)
}
