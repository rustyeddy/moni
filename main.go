package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/url"

	"github.com/gocolly/colly"
	"github.com/rustyeddy/store"
	log "github.com/sirupsen/logrus"
)

// Configuration manages all variables and parameters for a given run of moni.
type Configuration struct {
	Addrport   string
	ConfigFile string
	Changed    bool
	Daemon     bool
	Verbose    bool
}

var (
	config Configuration

	err     error
	acl     map[string]bool
	pages   map[url.URL]*Page
	storage *store.FileStore
)

func init() {
	flag.StringVar(&config.Addrport, "addr", "0.0.0.0:1212", "Address and port configuration")
	flag.StringVar(&config.ConfigFile, "config", "moni.json", "Moni config file")

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

	if storage, err = store.UseFileStore("."); err != nil || storage == nil {
		errFatal(err, "failed to useFileStore ")
	}

	// Process URLs on the command line if we have them
	if flag.Args() != nil {
		processURLs(flag.Args())
	}

	// Now run as a daemon if requested
	if !config.Daemon {
		var err error

		log.Println("listening on", config.Addrport)
		err = http.ListenAndServe(config.Addrport, nil)
		log.Fatal(err)
	}
}

func processURLs(urls []string) {
	// walk the command line arguments treating them as URLs
	for _, baseURL := range urls {

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
	pages[*r.URL] = NewPage(r.URL)
	fmt.Println("Request ", r.URL)
}

// Called after the response is
func doResponse(r *colly.Response) {
	fmt.Println("Response ", r.Request.URL)
}

func doScraped(r *colly.Response) {
	fmt.Printf("Scrap complete")
}
