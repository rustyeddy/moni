package main

import (
	"fmt"
	"net/url"
	"time"
)

// Site is the structure that maintains all pages for this
// root site
type Site struct {
	BaseURL string
	Pages   map[string]*Page
}

// Page represents a single web page with associate statistics
type Page struct {
	*url.URL
	Start  time.Time
	Finish time.Time
}

func processPage(urlstr string) (p *Page, err error) {
	var u *url.URL
	if u, err = url.Parse(urlstr); err != nil {
		panic(err)
	}
	fmt.Printf("URL: %+v\n", u)
	return p, err
}

