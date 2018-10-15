package main

import (
	"net/url"
	"time"

	log "github.com/sirupsen/logrus"
)

// ===================================================================
type Page struct {
	URL string

	Content []byte
	Links   map[string]*Page
	Ignored map[string]int

	LastCrawled time.Time
	StatusCode  int

	Start  time.Time
	Finish time.Time

	*Site
}

// ********************************************************************
type Pagemap map[string]*Page

// GetPage will sanitize the url, either find or create the
// corresponding page structure.  If the URL is deep, we also
// find the corresponding site structure.
func GetPage(url *url.URL) (pi *Page) {
	ustr := url.String()
	var ex bool
	if pi, ex = Pages[ustr]; !ex {
		pi = &Page{
			URL:     url.String(),
			Links:   make(map[string]*Page),
			Ignored: make(map[string]int),
			Site:    SiteFromURL(url),
		}
		Pages[ustr] = pi
		log.Infof("Created page %v", pi.URL)
	}
	return pi
}

func (pm Pagemap) Get(url string) (p *Page) {
	if p, e := Visited[url]; e {
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
	Visited[url] = p
}
