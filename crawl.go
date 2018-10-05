package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/gorilla/mux"

	//"github.com/gorilla/mux"
	"github.com/rustyeddy/store"
	log "github.com/sirupsen/logrus"
)

var (
	Visited    PageMap = make(PageMap)
	CrawlDepth int     = 1
	Store      *store.Store
)

func init() {
	var err error
	Store, err = store.UseStore("/srv/inv/")
	if err != nil {
		log.Fatalln("Could not use /srv/inv for storage ")
	}
}

// CrawlHandler will handle incoming HTTP request to crawl a URL
func CrawlHandler(w http.ResponseWriter, r *http.Request) {

	// 1. Extract and sanitize the URL from the request and sanitize it
	vars := mux.Vars(r)
	url := vars["url"]
	if !strings.HasPrefix("http", url) {
		url = "http://" + url
	}

	log.Infoln("crawl", url)
	page, err := Crawl(url)
	if err != nil {
		fmt.Fprintf(w, "url", err)
		return
	}

	// record the page now
	_, err = Store.StoreObject("page", page)
	if err != nil {
		log.Errorln("Failed to create local store")
	}

	jbytes, err := json.Marshal(page)
	if err != nil {
		log.Errorln("marshal json", url, err)
	}
	w.Header().Add("Content-Type", "application/json")
	w.Write(jbytes)
}

// Crawl will visit the given URL, and depending on configuration
// options potentially walk internal links.
func Crawl(url string) (p *Page, err error) {
	log.Infoln("crawling", url)

	// Create the collector and go get shit!
	c := colly.NewCollector(
		colly.MaxDepth(2),
	)

	c.OnRequest(func(r *colly.Request) {
		log.Infoln("  visiting site ", r.URL)
	})

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Request.AbsoluteURL(e.Attr("href"))
		if link == "" {
			return
		}
		p.Links[link]++
		e.Request.Visit(link)
	})

	c.OnResponse(func(r *colly.Response) {
		link := r.Request.URL
		log.Infoln("response recieved", link, r.StatusCode)
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
	if p = Visited.Get(url); p == nil {
		p = &Page{
			URL:   url,
			Links: make(map[string]int),
		}
		Visited[url] = p
	}
	p.Start = time.Now()
	c.Visit(url)
	return p, nil
}
