package main

import (
	"fmt"
	"net/url"
	"time"

	log "github.com/sirupsen/logrus"
)

// Page represents a single web page
type Page struct {
	url.URL    `json:"url"`
	StatusCode int                 `json:"statusCode"`
	Links      map[string][]string `json:"links"`

	ReqTime  time.Time     `json:"request"`
	RespTime time.Time     `json:"response"`
	Elapsed  time.Duration `json:"elapsed"`
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

func (p *Page) Print() {
	fmt.Printf("%s\n", p.URL.String())
	for l, _ := range p.Links {
		fmt.Printf("\t%s\n", l)
	}
	fmt.Printf("  elapsed time %v\n", p.Elapsed)
}
