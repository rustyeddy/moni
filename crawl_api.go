package moni

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

// registerCrawlers will register all handlers related to the
// crawling activities
func registerCrawlers(r *mux.Router) {
	r.HandleFunc("/crawl", CrawlListHandler) // Display "recent" crawl jobs
	r.HandleFunc("/crawl/", CrawlListHandler)
	r.HandleFunc("/crawl/{url}", CrawlUrlHandler).Methods("PUT", "POST", "GET", "DELETE") // Display a specific crawl job

}

// ServiceHandlers
// ========================================================================

// CrawlListHandler will return a list of all recent crawls.
// As stored in our storage (json) file.
func CrawlListHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, GetCrawls())
}

// CrawlHandler will handle incoming HTTP request to crawl a URL
func CrawlUrlHandler(w http.ResponseWriter, r *http.Request) {
	ustr := urlFromRequest(r)

	switch r.Method {

	case "GET":
		// Get crawl info about our site
		crawls := FindCrawls("crawl-*-" + ustr + "*")
		writeJSON(w, crawls)
		return

	case "PUT", "POST":
		// Create a crawl request
		// urlQ.Send(ustr)
		writeJSON(w, ustr)
		return

	case "DELETE":
		// Delete a site from this Crawler
		if _, ex := sites[ustr]; ex {
			delete(sites, ustr)
		}

	default:
		log.Errorf("CrawlUrlHandler unexpected method (%s) ", r.Method)
	}
	JSONError(w, fmt.Errorf("nothing good happened"))
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
