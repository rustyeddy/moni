package moni

import (
	"fmt"
	"net/http"
)

// REST API
// ====================================================================

func SiteListHandler(w http.ResponseWriter, r *http.Request) {
	if Sites == nil || len(Sites) < 1 {
		fmt.Fprintf(w, "no Sites ")
	}
	writeJSON(w, Sites)
}

func SiteIdHandler(w http.ResponseWriter, r *http.Request) {
	url := urlFromRequest(r)
	switch r.Method {
	case "GET":
		if site := Sites.Get(url); site != nil {
			writeJSON(w, site)
		} else {
			JSONError(w, ErrorNotFound(url+" site not found "))
		}

	case "PUT", "POST":
		AddNewSite(url)

	case "DELETE":
		Sites.Delete(url)

	default:
		JSONError(w, ErrorNotSupported("unspported method "+r.Method))
	}
}
