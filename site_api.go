package moni

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
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
	if Sites == nil || len(Sites) < 1 {
		fmt.Fprintf(w, "no Sites ")
	}
	writeJSON(w, Sites)
}

// SiteIdHandler manages requests targeted for a specific site.
func SiteIdHandler(w http.ResponseWriter, r *http.Request) {

	url := urlFromRequest(r)
	log.Debugln("SiteIdHandler request ", url)

	switch r.Method {
	case "GET":
		if site := Sites.Get(url); site != nil {
			writeJSON(w, site)
		} else {
			JSONError(w, ErrorNotFound(url+" site not found "))
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
		JSONError(w, ErrorNotSupported("unspported method "+r.Method))
	}
	return
}
