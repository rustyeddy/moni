package main

import (
	"net/url"
	"time"

	log "github.com/sirupsen/logrus"
)

var (
	pages map[string]*Page
)

// Page represents a single web page
type Page struct {
	*url.URL
	StatusCode int
	Links      map[string]int

	reqtime  time.Time
	resptime time.Time
}

// NewPage returns a new page structure
func newPage(u *url.URL) (p *Page) {

	p = &Page{
		URL:        u,
		StatusCode: 0,
		Links:      make(map[string]int),
	}
	return p
}

<<<<<<< HEAD
// NewPage returns a page structure that will hold all our cool stuff
func NewPage(u *url.URL) (p *Page) {
	p = &Page{
		URL: u,
	}
	log.Infof("New Page: %+v\n", u)
	return p
=======
// GetPage will return the page if it exists, or create otherwise.
func GetPage(urlstr string) (p *Page) {
	var ex bool

	u, err := normURL(urlstr)
	errPanic(err)

	if p, ex = pages[u.String()]; ex {
		return p
	}
	p = newPage(u)
	return p
}

// normURL is responsible for normalizaing the URL
func normURL(urlstr string) (u *url.URL, err error) {
	u, err = url.Parse(urlstr)
	return u, err
}

// errPanic something went wrong, panic.
func errPanic(err error) {
	if err != nil {
		panic(err)
	}
}

// nilPanic does so when it's parameter is such.
func nilPanic(val interface{}) {
	if val == nil {
		fmt.Printf("val is nil")
	}
>>>>>>> 2c23252715255f8758af33fe7b8b054f831f6d7d
}
