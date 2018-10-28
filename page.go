package moni

import (
	"fmt"
	"time"
)

const (
	NotCrawled = iota
	CrawlReady
	CrawlRequestSent
	CrawlResponseRecieved
	CrawlComplete
	CrawlErrored
	CrawlNotAllowed
)

// ===================================================================
type Page struct {
	URL string

	Content []byte
	Links   map[string]*Page
	Ignored map[string]int

	CrawlState int

	StatusCode int
	Err        error

	LastCrawled time.Time
	Start       time.Time
	Finish      time.Time
}

// Pagemap
// ********************************************************************
type Pagemap map[string]*Page

var (
	pages Pagemap
)

// String will represent the Page
// ====================================================================
func (p *Page) String() string {
	str := fmt.Sprintf("%s: lastcrawled: %s,  duration: %v links: %d ignored: %d\n", p.URL, p.LastCrawled, p.Finish, len(p.Links), len(p.Ignored))
	return str
}

func FetchPage(url string) *Page {
	panic("todo implement GetPage")
}

func StorePage(p *Page) {
	panic("todo implement save page")
}

// GetPage will sanitize the url, either find or create the
// corresponding page structure.  If the URL is deep, we also
// find the corresponding site structure.
func PageFromURL(ustr string) (pi *Page) {
	var ex bool
	if pi, ex = pages[ustr]; !ex {
		pi = &Page{
			URL:        ustr,
			Links:      make(map[string]*Page),
			Ignored:    make(map[string]int),
			CrawlState: CrawlReady,
		}
		if pages == nil {
			pages = make(map[string]*Page)
		}
		pages[ustr] = pi
	}
	return pi
}

func (pm Pagemap) Get(url string) (p *Page) {
	if p, e := pm[url]; e {
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
	pm[url] = p
}
