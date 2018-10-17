package main

import (
	"fmt"
	"time"
)

// ===================================================================
type Page struct {
	URL string

	Content []byte
	Links   map[string]*Page
	Ignored map[string]int

	StatusCode int

	LastCrawled time.Time
	Start       time.Time
	Finish      time.Time
	crawl       bool
	Err         error
}

// String will represent the Page
// ====================================================================
func (p *Page) String() string {
	str := fmt.Sprintf("%s: lastcrawled: %s,  duration: %v links: %d ignored: %d\n", p.URL, p.LastCrawled, p.Finish, len(p.Links), len(p.Ignored))
	return str
}

// ********************************************************************
type Pagemap map[string]*Page

// GetPage will sanitize the url, either find or create the
// corresponding page structure.  If the URL is deep, we also
// find the corresponding site structure.
func GetPage(ustr string) (pi *Page) {
	var ex bool
	if pi, ex = Pages[ustr]; !ex {
		pi = &Page{
			URL:     ustr,
			Links:   make(map[string]*Page),
			Ignored: make(map[string]int),
			crawl:   true,
		}
		Pages[ustr] = pi
	}
	return pi
}

func (pm Pagemap) Get(url string) (p *Page) {
	if p, e := Pages[url]; e {
		return p
	}
	return nil
}

func (pm Pagemap) Exists(url string) bool {
	if p := pm.Get(url); p != nil {
		return true
	}
	return false
}

func (pm Pagemap) Set(url string, p *Page) {
	Pages[url] = p
}
