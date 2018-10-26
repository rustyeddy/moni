package moni

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func processURL(ustr string) *Page {

	var page *Page

	// Normalize the URL and fill in a scheme it does not exist
	log.Infoln("processURL ", ustr)
	ustr, err := NormalizeURL(ustr)
	if err != nil {
		log.Infoln("processURL Normalize failed ", ustr)
		return nil
	}

	// This conversion back to string is necessary and simple domain
	// name like "example.com" will be placeded in the url.URL.Path
	// field instead of the Host field.  However url.String() makes
	// everything right.
	accessList.AllowHost(ustr)

	// If the url has too recently been scanned we will return
	// null for the job, however a copy of the scan will is
	// available and will be returned to the caller.
	if page = CrawlOrNot(ustr); page == nil {
		log.Infoln("processURL rejected url ", ustr)
		return nil
	}

	log.Infoln("returning from process URL with page ", page.URL)
	return page
}

// NormalizeURL will make sure a scheme (protocol) is prepended
// go domains that do not already sport a scheme.
func NormalizeURL(urlstr string) (ustr string, err error) {
	if urlstr == "" {
		return "", errors.New("blank URL")
	}

	// Parse the string into the parts of the url
	u, err := url.Parse(urlstr)
	if err != nil {
		return "", fmt.Errorf("parse url %s -> %v", urlstr, err)
	}

	// make sure only have either http and https schemes. If the url has
	// no scheme, colly will reject it, so we'll add the http scheme.
	switch u.Scheme {
	case "http", "https":
		// we are good
	case "":
		u.Scheme = "http"
	default:
		accessList.Unsupported[urlstr]++
		return "", fmt.Errorf("Error Not Supported")
	}
	ustr = u.String()
	return ustr, nil
}

func urlFromRequest(r *http.Request) string {
	// Extract the url(s) that we are going to walk
	vars := mux.Vars(r)
	ustr := vars["url"]

	// Normalize the URL and fill in a scheme it does not exist
	ustr, err := NormalizeURL(ustr)
	if err != nil {
		log.Errorf("I had a problem with the url %v", ustr)
		return ""
	}
	return ustr
}

// GetHostname from a URL type string
func GetHostname(h string) (host string) {
	var err error
	var u *url.URL

	// why does url.Parse move the hostname 'gum.com' to a page?
	if u, err = url.Parse(h); err != nil {
		log.Errorln("failed to parse hostname", err)
		return ""
	}
	// TODO Ugly stuff see above
	if u.Host == "" && u.Scheme == "" && u.Path != "" {
		u.Host = u.Path
		u.Path = ""
	}
	return u.Hostname()
}
