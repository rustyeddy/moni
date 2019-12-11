package main

import (
	"net/url"
	"time"

	log "github.com/sirupsen/logrus"
)

// Page represents a single web page
type Page struct {
	url.URL
	StatusCode int
	Links      map[string][]string

	ReqTime  time.Time
	RespTime time.Time
	Elapsed  time.Duration
}

// NewPage returns a page structure that will hold all our cool stuff
func NewPage(u url.URL) (p *Page) {
	p = &Page{
		URL:   u,
		Links: make(map[string][]string),
	}
	log.Infof("New Page: %+v for a total of %d\n", u, len(pages))
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
