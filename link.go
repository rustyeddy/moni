package main

import (
	"net/url"
	"time"

	log "github.com/sirupsen/logrus"
)

type Link struct {
	URL        string        `json:"url"`
	Anchors    []string      `json:"anchors"`
	Internal   bool          `json:"internal"`
	Reachable  bool          `json:"reachable"`
	Response   time.Duration `json:"response"`
	StatusCode int           `json:"statuscode"`
	Last       time.Time     `json:"last"`

	*Page `json:"-"`
}

// NewLink creates a new link that will be collected by the parent page
func NewLink(urlstr string) (l *Link) {
	var u *url.URL

	if u = scrubURL(urlstr); u == nil {
		log.Infof("NewLink: URL has been removed %s ")
		return nil
	}
	l = &Link{
		URL: u.String(),
	}
	return l
}
