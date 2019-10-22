package main

import (
<<<<<<< HEAD
	"flag"
	"fmt"
	"net/url"

	"github.com/gocolly/colly"
	"github.com/rustyeddy/store"
	log "github.com/sirupsen/logrus"
)

// Configuration manages all variables and parameters for a given run of moni.
type Configuration struct {
	ConfigFile string
	Changed    bool
	Verbose    bool
}
=======
	"net/http"
>>>>>>> 2c23252715255f8758af33fe7b8b054f831f6d7d

	log "github.com/sirupsen/logrus"
)

func init() {
<<<<<<< HEAD
	var err error

	flag.StringVar(&config.ConfigFile, "config", "moni.json", "Moni config file")

	storage, err = store.UseFileStore(".")
	errPanic(err)

	pages = make(map[string]*Page)
	acl = make(map[string]bool)

	// TODO read the acls from a file
	acl["localhost"] = false
}

func main() {
	flag.Parse()

	urls := flag.Args()
	if urls == nil || len(urls) == 0 {
		log.Fatal("Expected some sites, got none")
	}
	processURLs(urls)
}

func processURLs(urls []string) {

	// walk the command line arguments treating them as URLs
	for _, baseURL = range urls {

		// Place the command line url in the acl allowed list
		if config.Verbose {
			log.Infof("Add website %s to ACL\n", baseURL)
		}

		u, err := url.Parse(baseURL)
		errPanic(err)

		// This is a little risky
		acl[u.Hostname()] = true

		// This will become sending a message
		Walk(u)
	}
	storage.Save("config2.json", config)
	fmt.Println("The end...")
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
	}
	fmt.Printf("Hostname %s ", host)
	if ok, ex = acl[host]; ex {
		return ok
	}
	return false
}

// Crawl the given URL
func Walk(u *url.URL) {
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
	if u == nil {
		return
	}

	fmt.Printf("url: %+v ", u)
	if processPage(urlstr) {
		fmt.Println("ok ...")
		e.Request.Visit(urlstr)
	} else {
		fmt.Println(" blocked ...")
	}
}

// Called before the request is sent
func doRequest(r *colly.Request) {
	pages[r.URL.String()] = NewPage(r.URL)
	fmt.Println("Request ", r.URL)
}

// Called after the response is
func doResponse(r *colly.Response) {
	fmt.Println("Response ", r.Request.URL)
}

func doScraped(r *colly.Response) {
	fmt.Println("Scraped ", r.Request.URL)
}
=======
	log.SetFormatter(&log.JSONFormatter{})
}

func main() {
	// example usage: curl -s 'http://127.0.0.1:7171/?url=http://go-colly.org/'
	addr := ":7171"

	http.HandleFunc("/", handler)
>>>>>>> 2c23252715255f8758af33fe7b8b054f831f6d7d

	log.Println("listening on", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func errPanic(err error) {
	if err != nil {
		panic(err)
	}
}
