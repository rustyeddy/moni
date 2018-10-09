package main

import (
	"net/url"
	"strings"

	log "github.com/sirupsen/logrus"
)

// PreparseURL will clean up the URL if needed and determine
// if the page represented by the URL needs scanning.  For example,
// if the URL has just been crawled, we do not want to crawl it again.
// Or it could be a URL that we do not want to crawl (like amazon.com for
// example).
func PrepareURL(urlstr string) (pi *PageInfo) {

	// Parse the string into the parts of the url
	u, err := url.Parse(urlstr)
	if err != nil {
		log.Errorln("failed to parse url", urlstr, err)
		return
	}
	log.Debugf("PrepareURL URL: %+v", u)

	// Check if we will allow crawling this hostname
	if ACL.Reject(u.Hostname()) {
		log.Infoln("ACL Rejects hostname", urlstr)
		return
	}

	// Get (or create) a PageInfo struct for the given URL
	pi = Pages.Get(u.Hostname())
	if pi == nil {
		log.Errorln("Expected page for %s got ()", urlstr)
		return
	}

	// If the page has recently been crawled we wont do it again
	if pi.Crawled {
		log.Debugln("rejecting url", urlstr)
		return
	}

	// make sure only have either http and https schemes. If the url has
	// no scheme, colly will reject it, so we'll add the http scheme.
	switch u.Scheme {
	case "":
		u.Scheme = "http"
	case "http", "https":
	default:
		log.Warn("unsupported scheme", u.Scheme)
		return
	}

	return pi
}

func nameFromURL(urlstr string) (name string) {
	u, err := url.Parse(urlstr)
	if err != nil {
		log.Errorln(err)
		return
	}
	name = u.Hostname()
	name = strings.Replace(name, ".", "-", -1)
	return name
}
