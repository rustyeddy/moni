package main

import (
	"net/url"
	"strings"

	log "github.com/sirupsen/logrus"
)

// NormalizeURL will make sure a scheme (protocol) is prepended
// go domains that do not already sport a scheme.
func NormalizeURL(urlstr string) (u *url.URL, ustr string) {
	var err error

	// Parse the string into the parts of the url
	u, err = url.Parse(urlstr)
	if err != nil {
		log.Errorln("reject url ~ failed parse", urlstr, err)
		return nil, ""
	}
	log.Debugf("PrepareURL URL: %+v, hostname %s", u, u.Hostname())

	// make sure only have either http and https schemes. If the url has
	// no scheme, colly will reject it, so we'll add the http scheme.
	switch u.Scheme {
	case "":
		u.Scheme = "http"
	case "http", "https":
	default:
		ACL.Unsupported[urlstr]++
		log.Infof("unsupported scheme %s", u.Scheme)
		return nil, ""
	}
	return u, u.String()
}

// PreparseURL will clean up the URL if needed and determine
// if the page represented by the URL needs scanning.  For example,
// if the URL has just been crawled, we do not want to crawl it again.
// Or it could be a URL that we do not want to crawl (like amazon.com for
// example).
func PrepareURL(urlstr string) (pi *Page) {

	u, ustr := NormalizeURL(urlstr)
	if u == nil {
		log.Errorln("failed to normalize ", urlstr)
	}

	// Check if we will allow crawling this hostname
	allowed := ACL.IsAllowed(ustr)
	if !allowed {
		log.Infof("reject url %s not allowed", ustr)
		return nil
	}
	// Get (or create) a PageInfo struct for the given URL
	pi = GetPage(ustr)
	if pi == nil {
		log.Errorln("rejected url ~ get page failed", ustr)
		return nil
	}

	// If the page has recently been crawled we wont do it again
	if pi.Crawled {
		log.Debugln("rejecting url ~ recently crawled", ustr)
		return nil
	}

	// Ugly should be more better
	pi.URL = ustr
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
