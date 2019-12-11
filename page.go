package main

import (
	"encoding/json"
	"io"
	"net/url"
	"os"
	"time"

	"github.com/gocolly/colly"
	log "github.com/sirupsen/logrus"
)

// Page represents a single web page
type Page struct {
	url.URL
	StatusCode int
	Links      map[string]int

	ReqTime  time.Time
	RespTime time.Time
	Elapsed  time.Duration
}

// NewPage returns a page structure that will hold all our cool stuff
func NewPage(u url.URL) (p *Page) {
	p = &Page{
		URL:   u,
		Links: make(map[string]int),
	}
	log.Infof("New Page: %+v\n", u)
	pages[u] = p
	return p
}

// GetPage will return the page if it exists, or create otherwise.
func GetPage(u url.URL) (p *Page) {
	var ex bool
	if p, ex = pages[u]; ex {
		return p
	}
	p = NewPage(u)
	return p
}

// Crawl the given URL
func (p *Page) Walk(w io.Writer) {
	var urls []string

	c := colly.NewCollector()

	log.Infof("Visiting page %s", p.URL.String())

	// Setup all the collbacks
	c.OnHTML("a", func(e *colly.HTMLElement) {
		refurl := e.Attr("href")
		link := e.Request.AbsoluteURL(refurl)
		p.Links[link]++
	})

	c.OnRequest(func(r *colly.Request) {
		pages[*r.URL] = NewPage(*r.URL)
		log.Infoln("Request ", r.URL)
	})

	c.OnResponse(func(r *colly.Response) {
		p.RespTime = time.Now()
		log.Infoln("Response ", r.Request.URL)
	})

	c.OnScraped(func(r *colly.Response) {
		log.Debugf("\tLinks: %s", p.URL.String())
		for ustr, _ := range p.Links {
			log.Debugf("\t~> %s", ustr)
			if config.Recurse {
				urls = append(urls, ustr)
			}
		}
	})

	p.ReqTime = time.Now()
	c.Visit(p.String())
	p.RespTime = time.Now()
	p.Elapsed = p.RespTime.Sub(p.ReqTime)

	if w == nil {
		w = os.Stdout
	}

	links := []string{}
	for l, _ := range p.Links {
		links = append(links, l)
	}

	resp := struct {
		Status     string        `json:"status"`
		ReturnCode int           `json:"returnCode"`
		URL        string        `json:"url"`
		Links      []string      `json:"links"`
		Elapsed    time.Duration `json:"elapsed"`
	}{
		Status:     "",
		ReturnCode: 200,
		URL:        p.URL.String(),
		Links:      links,
		Elapsed:    p.Elapsed,
	}

	json.NewEncoder(w).Encode(resp)
}
