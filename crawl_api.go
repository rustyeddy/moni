package moni

import (
	"errors"
	"net/http"

	"github.com/gorilla/mux"
)

// registerCrawlers will register all handlers related to the
// crawling activities
func registerCrawlers(r *mux.Router) {
	r.HandleFunc("/crawlids", CrawlListHandler)    // Display "recent" crawl jobs
	r.HandleFunc("/crawl/{url}", CrawlHandler)     // Create a (recurring) crawl job for url
	r.HandleFunc("/crawlid/{cid}", CrawlIdHandler) // Display a specific crawl job
}

// ServiceHandlers
// ========================================================================

// CrawlHandler will handle incoming HTTP request to crawl a URL
func CrawlHandler(w http.ResponseWriter, r *http.Request) {
	ustr := urlFromRequest(r)
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
	vars := mux.Vars(r)
	cid := vars["cid"]

	var page *Page
	if page := FetchPage(cid); page != nil {
		JSONError(w, errors.New("page fetch failed for cid "+cid))
		return
	}
	writeJSON(w, page)
}
