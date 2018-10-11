package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gocolly/colly"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type CrawlJob struct {
	ID  string
	URL string
	*Page
}

var (
	Visited    Pagemap = make(Pagemap)
	CrawlDepth int     = 1
)

// CrawlHandler will handle incoming HTTP request to crawl a URL
func CrawlHandler(w http.ResponseWriter, r *http.Request) {

	// 1. Extract and sanitize the URL from the request and sanitize it
	vars := mux.Vars(r)
	_, ustr := NormalizeURL(vars["url"])

	// Add this URL to the allowed list so PrepareURL does not reject
	// this url
	ACL.AllowHost(ustr)

	// Now prepare the URL and get a page back
	pi := PrepareURL(ustr)
	if pi == nil {
		log.Errorf("url rejected %s", ustr)
		fmt.Fprintf(w, "url rejected %s", ustr)
		return
	}

	// TODO: this is where we create a Job and give them a token
	// to look back later.
	log.Infoln("crawl", ustr)
	page, err := Crawl(ustr)
	if err != nil {
		fmt.Fprint(w, "ustr", err)
		return
	}

	// Determine an index to store the page under. The URL is perfect
	// except that it will likely contain '/' which conflict with the
	// pathname.  Hence our index must not contain slashes.
	name := nameFromURL(ustr)
	_, err = Storage.StoreObject(name, page)
	if err != nil {
		log.Errorln("Failed to create local store")
		fmt.Fprintf(w, "Internal Error %s", ustr)
		return
	}

	jbytes, err := json.Marshal(page)
	if err != nil {
		log.Errorln("marshal json", ustr, err)
		fmt.Fprintf(w, "Internal Error")
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(jbytes)
}

// Crawl will visit the given URL, and depending on configuration
// options potentially walk internal links.
func Crawl(urlstr string) (p *Page, err error) {
	log.Infoln("crawling", urlstr)

	// Create the collector and go get shit!
	c := colly.NewCollector(
		colly.MaxDepth(2),
	)

	c.OnRequest(func(r *colly.Request) {
		log.Infoln("Visiting site ", r.URL)
	})

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Request.AbsoluteURL(e.Attr("href"))
		if link == "" {
			return
		}
		pi := PrepareURL(link)
		if pi == nil {
			// The link has been filtered for one
			// reason or another, we will move along
			p.Ignored[link]++
			return
		}
		p.Links[link] = GetPage(link)
		e.Request.Visit(link)
	})

	c.OnResponse(func(r *colly.Response) {
		link := r.Request.URL
		log.Infoln("  response from", link, "status", r.StatusCode)
		p.StatusCode = r.StatusCode
		p.End = time.Now()
		p.Crawled = true
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Infoln("error:", r.StatusCode, err)
		p.StatusCode = r.StatusCode
		p.End = time.Now()
	})

	// Get the page we'll use for this walk
	if p = Visited.Get(urlstr); p == nil {
		p = &Page{
			URL:     urlstr,
			Links:   make(map[string]*Page),
			Ignored: make(map[string]int),
		}
		Visited[urlstr] = p
	}
	p.Start = time.Now()
	c.Visit(urlstr)
	return p, nil
}
