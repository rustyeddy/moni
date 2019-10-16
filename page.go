package main

import (
	"fmt"
	"net/url"
	"time"
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
}
