package main

import (
	"net/url"

	log "github.com/sirupsen/logrus"
)

type Link struct {
	*url.URL          // URL of this link
	*Page             // Option link to the actual page
	Anchors  []string // string of anchor text used to point to this link
	Internal bool     // Is this link internal or external
}

// NewLink creates a new link that will be collected by the parent page
func NewLink(urlstr string) (l *Link) {
	var u *url.URL

	if u = scrubURL(urlstr); u == nil {
		log.Infof("NewLink: URL has been removed %s ")
		return nil
	}
	l = &Link{
		URL: u,
	}
	return l
}
