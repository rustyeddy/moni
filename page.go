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
	TimeStamp  `json:"timestamp"`
	TimeStamps []TimeStamp `json:"timestamps"`
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

func (p *Page) PageString() (s string) {
	s = fmt.Sprintf("%-40s: links: %-4d resp: %-10v ", p.URL.String(), len(p.Links), p.Elapsed)
	if tslen := len(p.TimeStamps); tslen > 0 {
		if tslen > 4 {
			tslen -= 4
		}
		s += fmt.Sprintf("\tlast elasped: %v", p.TimeStamps[tslen:])
	}
	s += "\n"

	if config.Verbose {
		for l, _ := range p.Links {
			s += fmt.Sprintf("\t%s\n", l)
		}
	}
	return s
}
