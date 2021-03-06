package main

import (
	"fmt"
	"net/url"
	"time"
)

// Pages is a collection (map) of pages that belong to the
// same website.
type Pages map[url.URL]*Page

// Page represents a single web page
type Page struct {
	*Site      `json:"-"`
	url.URL    `json:"url"`
	Links      map[string]*Link `json:"links"`
	StatusCode int              `json:"statusCode"`

	TimeStamp  `json:"timestamp"`
	TimeStamps []TimeStamp `json:"timestamps"`
	WalkTimer  *time.Timer
}

type PageInfo struct {
	URL      string        `json:"url"`
	Links    []*Link       `json"links"`
	Response time.Duration `json:"duration"`
}

// NewPage will create a new page based on the URL, prepare the
// Links map.
func NewPage(u *url.URL) (p *Page) {
	p = &Page{
		URL:   *u,
		Links: make(map[string]*Link),
	}
	return p
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
