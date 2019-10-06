package main

import (
	"fmt"
	"net/url"
	"time"
)

type allPages map[string]*Page

// Page represents a single web page with associate statistics
type Page struct {
	*url.URL

	Start  time.Time
	Finish time.Time

	Ignore bool
}

// NewPage returns a page structure that will hold all our cool stuff
func NewPage(urlstr string) (p *Page) {
	var u *url.URL
	var err error

	if u, err = url.Parse(urlstr); err != nil {
		panic(err)
	}
	p = &Page{
		URL: u,
	}
	fmt.Printf("URL: %+v\n", u)
	return p
}
