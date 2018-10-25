package moni

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func urlWatcher(urlq chan string, crawlq chan *Page, errq chan error) {

	for {
		url := <-urlq

		ustr, err := NormalizeURL(url)
		if err != nil {
			errq <- fmt.Errorf("Error Normalizing url %s err %v", url, err)
			// Aint no crawling gonna happen around heree
			continue
		}

		// CrawlOrNot will return a *Page with Page.CrawlState = CrawlReady
		pg := CrawlOrNot(ustr)
		if pg == nil {
			// No crawling going on round here ...
			continue
		}

		// We are good to go, send this page off too the crawler
		crawlq <- pg
	}

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
