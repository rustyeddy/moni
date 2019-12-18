package main

import (
	"fmt"
	"net/url"

	log "github.com/sirupsen/logrus"
)

// Page represents a single web page
type Page struct {
	url.URL    `json:"url"`
	Links      map[string][]string `json:"links"`
	StatusCode int                 `json:"statusCode"`
	TimeStamp
	TimeStamps []TimeStamp
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

func (p *Page) String() string {
	return p.URL.String()
}

func (p *Page) Print() {
	fmt.Printf("%40s: resp: %v links: %d\n", p.URL.String(), p.Elapsed, len(p.Links))
	if config.Verbose {
		for l, _ := range p.Links {
			fmt.Printf("\t%s\n", l)
		}
	}
}
