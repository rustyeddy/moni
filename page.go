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
	CrawlReady bool
	StatusCode int

	Err error

	LastCrawled time.Time
	Start       time.Time
	Finish      time.Time
}

// Pagemap
// ********************************************************************
type Pagemap map[string]*Page

func initPages() (p Pagemap) {
	return make(Pagemap)
}

// NewPage returns a newly created page represented by the URL, NewPage
// registers itself the pages Pagemap.
func NewPage(url string) (p *Page) {
	p = &Page{
		URL:        url,
		Links:      make(map[string]*Page),
		Ignored:    make(map[string]int),
		CrawlReady: true,
	}
	pages[url] = p
	return p
}

func GetPage(url string) (p *Page) {
	if p = pages.Get(url); p == nil {
		p = NewPage(url)
	}
	return
}

func removeTrailingSlash(u string) string {
	l := len(u) - 1
	if u[l] == '/' {
		return u[:l-1]
	}
	return u
}

// FetchPage returns the page from the pagemap if it exists. If
// it does not exist, nil will be returned.
func FetchPage(url string) (p *Page) {
	return pages.Get(url)
}

// StorePage will save the page to the pagemap, if the page index does
// not exist it will be created for the page.  If the page already
// exists it will be overwritten with the new page.
func StorePage(p *Page) {
	pages[p.URL] = p
}

// GetPage will sanitize the url, either find or create the
// corresponding page structure.  If the URL is deep, we also
// find the corresponding site structure.
func PageFromURL(ustr string) (pi *Page) {
	var ex bool

	ustr = removeTrailingSlash(ustr)
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

// String will represent the Page
// ====================================================================
func (p *Page) String() string {
	str := fmt.Sprintf("%s: lastcrawled: %s,  duration: %v links: %d ignored: %d\n", p.URL, p.LastCrawled, p.Finish, len(p.Links), len(p.Ignored))
	return str
}

func (pm Pagemap) Get(url string) (p *Page) {
	url = removeTrailingSlash(url)
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
