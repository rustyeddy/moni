package main

import (
	"net/url"
	"time"

	log "github.com/sirupsen/logrus"
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
func NewPage(u *url.URL) (p *Page) {
	p = &Page{
		URL: u,
	}
	log.Infof("New Page: %+v\n", u)
	return p
}
