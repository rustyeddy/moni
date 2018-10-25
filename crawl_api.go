package moni

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

// ServiceHandlers
// ========================================================================

// CrawlHandler will handle incoming HTTP request to crawl a URL
func CrawlHandler(w http.ResponseWriter, r *http.Request) {

	// Prepare for Execution
	// Extract the url(s) that we are going to walk
	vars := mux.Vars(r)
	ustr := vars["url"]

	// Normalize the URL and fill in a scheme it does not exist
	ustr, err := NormalizeURL(ustr)
	if err != nil {
		fmt.Fprintf(w, "I had a problem with the url %v", ustr)
		return
	}

	// This conversion back to string is necessary and simple domain
	// name like "example.com" will be placeded in the url.URL.Path
	// field instead of the Host field.  However url.String() makes
	// everything right.
	accessList.AllowHost(ustr)

	// If the url has too recently been scanned we will return
	// null for the job, however a copy of the scan will is
	// available and will be returned to the caller.
	page := CrawlOrNot(ustr)
	if page == nil {
		log.Errorf("url rejected %s", ustr)
		fmt.Fprintf(w, "url rejected %s", ustr)
		return
	}

	Crawl(page)

	writeJSON(w, page)

	// Cache the results ...  We'll replace '/' with '-' and
	// store the results in the cache store.
	storePageCrawl(page)
}

// CrawlListHandler will return a list of all recent crawls.
// As stored in our storage (json) file.
func CrawlListHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, GetCrawls())
}

// CrawlIdHandler will return a list of all recent crawls.
// As stored in our storage (json) file.
func CrawlIdHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the url(s) that we are going to walk
	vars := mux.Vars(r)
	cid := vars["cid"]
	st := GetStorage()

	page := new(Page)
	_, err := st.FetchObject(cid, page)
	if err != nil {
		JSONError(w, err)
		return
	}
	writeJSON(w, page)
}
