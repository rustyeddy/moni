package moni

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// REST API
// ====================================================================

// Register the site routes
func registerSites(r *mux.Router) {
	r.HandleFunc("/sites", SiteListHandler).Methods("GET")
	r.HandleFunc("/site/{url}", SiteIdHandler).Methods("GET", "POST", "PUT", "DELETE")
}

// SiteListHandler may respond with multiple Site entries
func SiteListHandler(w http.ResponseWriter, r *http.Request) {
	if sites == nil || len(sites) < 1 {
		fmt.Fprintf(w, "")
	}
	writeJSON(w, sites)
}

// SiteIdHandler manages requests targeted for a specific site.
func SiteIdHandler(w http.ResponseWriter, r *http.Request) {

	url := urlFromRequest(r)
	log.Debugln("SiteIdHandler request ", url)

	switch r.Method {
	case "GET":
		if site := sites.Get(url); site != nil {
			writeJSON(w, site)
		} else {
			JSONError(w, fmt.Errorf("site not found %s", url))
		}

	case "PUT", "POST":
		if nsite := AddNewSite(url); nsite != nil {
			writeJSON(w, nsite)
		} else {
			JSONError(w, errors.New("Failed to Add "+url))
		}

	case "DELETE":
		h := GetHostname(url)
		RemoveSite(h)

	default:
		JSONError(w, fmt.Errorf("unspported method "+r.Method))
	}
	return
}
