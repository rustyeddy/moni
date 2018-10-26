package moni

import (
	"net/http"

	"github.com/gorilla/mux"
)

// registerCrawlers will register all handlers related to the
// crawling activities
func registerCrawlers(r *mux.Router) {
	r.HandleFunc("/acl", ACLHandler)               // Display ACLs
	r.HandleFunc("/crawlids", CrawlListHandler)    // Display "recent" crawl jobs
	r.HandleFunc("/crawl/{url}", CrawlHandler)     // Create a (recurring) crawl job for url
	r.HandleFunc("/crawlid/{cid}", CrawlIdHandler) // Display a specific crawl job
}

// ServiceHandlers
// ========================================================================

// CrawlHandler will handle incoming HTTP request to crawl a URL
func CrawlHandler(w http.ResponseWriter, r *http.Request) {

	// Prepare for Execution
	// Extract the url(s) that we are going to walk
	vars := mux.Vars(r)
	ustr := vars["url"]

	Crawler.UrlQ <- ustr
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

	page := new(Page)
	_, err := storage.FetchObject(cid, page)
	if err != nil {
		JSONError(w, err)
		return
	}
	writeJSON(w, page)
}
