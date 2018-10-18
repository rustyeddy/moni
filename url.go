package main

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

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
		// fallthrough
	case "":
		u.Scheme = "http"
	default:
		ACL.Unsupported[urlstr]++
		return "", ErrorNotSupported(urlstr)
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
