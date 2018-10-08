package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

var (
	Visited    PageMap = make(PageMap)
	CrawlDepth int     = 1
)

// CrawlHandler will handle incoming HTTP request to crawl a URL
func CrawlHandler(w http.ResponseWriter, r *http.Request) {

	// 1. Extract and sanitize the URL from the request and sanitize it
	vars := mux.Vars(r)
	url := vars["url"]
	if !strings.HasPrefix("http", url) {
		url = "http://" + url
	}

	log.Infoln("crawl", url)
	// TODO: this is where we create a Job and give them a token
	// to look back later.
	page, err := Crawl(url)
	if err != nil {
		fmt.Fprintf(w, "url", err)
		return
	}

	// Determine an index to store the page under. The URL is perfect
	// except that it will likely contain '/' which conflict with the
	// pathname.  Hence our index must not contain slashes.
	st := Storage()
	name := nameFromURL(url)
	_, err = st.StoreObject(name, page)
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
		u, err := url.Parse(link)
		if err != nil {
			log.Warnln("failed to parse link", link, err)
			return
		}

		switch u.Scheme {
		case "http", "https", "":
			p.Links[link]++
			e.Request.Visit(link)

		default:
			// Ignore these links and do not crawl them
			p.Ignored[link]++
			if p.Ignored[link] < 3 {
				log.Warnln("ignoring scheme ", u.Scheme)
			}
		}
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
			Links:   make(map[string]int),
			Ignored: make(map[string]int),
		}
		Visited[urlstr] = p
	}
	p.Start = time.Now()
	c.Visit(urlstr)
	return p, nil
}

func nameFromURL(urlstr string) (name string) {
	u, err := url.Parse(urlstr)
	if err != nil {
		log.Errorln(err)
		return
	}
	name = u.Hostname()
	name = strings.Replace(name, ".", "-", -1)
	return name
}
