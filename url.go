package moni

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

//			      URL Cleansing and Normalization
// ====================================================================

// helper extract url from string
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
		acl.Unsupported[urlstr]++ // via AccessList
		return "", fmt.Errorf("Error Not Supported")
	}
	ustr = u.String()
	return ustr, nil
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

func processURL(url string) (page *Page) {

	log.Infoln("processURL ", url)

	// normalize the URL
	urlstr, err := NormalizeURL(url)
	if err != nil {
		log.Errorf("url normaization failed %v", err)
		return nil
	}

	log.Infoln("  fetching page ", urlstr)
	page = FetchPage(urlstr)
	if page == nil {
		log.Infoln("    not found ~ create a new one ..")
		if page = NewPage(urlstr); page == nil {
			log.Infoln("      Still could not find page for ", urlstr, " failing ...")
			return nil
		}

		if page == nil {
			panic("PAGE IS NIL ~ NEVER!")
		}
	}

	log.Infoln("  ACL is this site allowed:  ", page.URL)
	if !acl.IsAllowed(page.URL) {
		log.Infof("  URL %s NOT allowed, skipping ... ", page.URL)
		return nil
	}
	log.Infoln(page.URL, " is ready for crawling... ")
	return page
}
