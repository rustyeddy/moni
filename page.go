package main

import (
	"fmt"
	"net/url"
	"time"

	"github.com/gocolly/colly"
	log "github.com/sirupsen/logrus"
)

// Page represents a single web page
type Page struct {
	*url.URL
	StatusCode int
	Links      map[string]int

	ReqTime  time.Time
	RespTime time.Time
}

// NewPage returns a page structure that will hold all our cool stuff
func NewPage(u *url.URL) (p *Page) {
	p = &Page{
		URL:   u,
		Links: make(map[string]int),
	}
	log.Infof("New Page: %+v\n", u)
	pages[*u] = p
	return p
}

// GetPage will return the page if it exists, or create otherwise.
func GetPage(urlstr string) (p *Page) {
	var ex bool

	u, err := url.Parse(urlstr)
	errPanic(err)

	if p, ex = pages[*u]; ex {
		return p
	}
	p = NewPage(u)
	return p
}

// Crawl the given URL
func (p *Page) Walk() {
	c := colly.NewCollector()

	// Setup all the collbacks
	c.OnHTML("a", func(e *colly.HTMLElement) {
		refurl := e.Attr("href")
		link := e.Request.AbsoluteURL(refurl)
		p.Links[link]++
		//p.Links[refurl]++
		fmt.Println("link: ", link)
	})

	c.OnRequest(func(r *colly.Request) {
		pages[*r.URL] = NewPage(r.URL)
		fmt.Println("Request ", r.URL)
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Response ", r.Request.URL)
	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Printf("page and links: %s\n", p.URL.String())
		for ustr, _ := range p.Links {
			fmt.Printf("\t~> %s\n", ustr)
			if config.Recurse {
				c.Visit(ustr)
			}
		}
	})

	c.Visit(p.URL.String())
}
