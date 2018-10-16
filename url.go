package main

import (
	"errors"
	"fmt"
	"net/url"
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
