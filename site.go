package main

import (
	"fmt"
	"github.com/gocolly/colly"	
)

var (
	AllSites map[string]Site
)

// Site is the structure that maintains all pages for this
// root site
type Site struct {
	BaseURL string
	Pages   map[string]*Page

	Crawled bool
}

// NewSite will create a new Site structure 
func NewSite(baseurl string) (s *Site) {
	s = &Site{
		BaseURL: baseurl,
		Pages: make(map[string]*Page),

		Crawled: false,
	}
	return s
}

// Crawl the given URL
func (s *Site) Crawl() {
	c := colly.NewCollector()

	// Setup all the collbacks
	c.OnHTML("a", doHyperlink)
	c.OnRequest(doRequest)

	fmt.Println("Visiting: ", s.BaseURL)
	c.Visit(s.BaseURL)
}

func doHyperlink(e *colly.HTMLElement) {
	thing := e.Attr("href")
	page, err := processPage(thing)
	err_panic(err)
	if page != nil {
		e.Request.Visit(thing)
	}
}

func doRequest(r *colly.Request) {
	fmt.Println("Request ", r.URL)
}
