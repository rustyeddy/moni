package moni

import (
	"net/http"

	"github.com/gorilla/mux"
)

// ServiceHandlers
// ========================================================================

// CrawlHandler will handle incoming HTTP request to crawl a URL
func CrawlHandler(w http.ResponseWriter, r *http.Request) {

	// Prepare for Execution
	// Extract the url(s) that we are going to walk
	vars := mux.Vars(r)
	ustr := vars["url"]

	sched.URLQ <- ustr
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
