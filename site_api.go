package moni

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

// REST API
// ====================================================================

// Register the site routes
func registerSites(r *mux.Router) {
	r.HandleFunc("/site", HostSiteHandler).Methods("GET")
	r.HandleFunc("/site/", SiteListHandler).Methods("GET")
	r.HandleFunc("/site/{url}", SiteIdHandler).Methods("GET", "POST", "PUT", "DELETE")
}

// SiteListHandler may respond with multiple Site entries
func SiteListHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, FetchSites())
}

// SiteIdHandler manages requests targeted for a specific site.
func SiteIdHandler(w http.ResponseWriter, r *http.Request) {

	url := urlFromRequest(r)
	log.Debugln("SiteIdHandler request ", url)

	switch r.Method {
	case "GET":
		if site := FetchSite(url); site != nil {
			writeJSON(w, site)
		} else {
			JSONError(w, fmt.Errorf("site not found %s", url))
		}

	case "PUT", "POST":
		// Need a little more fan fair
		Crawler.UrlQ <- url
		writeJSON(w, map[string]string{"saved": url})

	case "DELETE":
		h := GetHostname(url)
		DeleteSite(h)

	default:
		JSONError(w, fmt.Errorf("unspported method "+r.Method))
	}
	return
}

// HostSiteHandler returns data about this particular site (server really)
func HostSiteHandler(w http.ResponseWriter, r *http.Request) {

}
