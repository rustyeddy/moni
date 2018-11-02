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
	r.HandleFunc("/site", SiteListHandler).Methods("GET")
	r.HandleFunc("/site/", SiteListHandler).Methods("GET")
	r.HandleFunc("/site/{url}", SiteIdHandler).Methods("GET", "POST", "PUT", "DELETE")
}

// SiteListHandler may respond with multiple Site entries
func SiteListHandler(w http.ResponseWriter, r *http.Request) {
	log.Infof("SiteListHandler sites %+v", sites)
	if len(sites) < 1 {
		sites = Sitemap{}
	}
	writeJSON(w, sites)
}

// SiteIdHandler manages requests targeted for a specific site.
func SiteIdHandler(w http.ResponseWriter, r *http.Request) {
	url := urlFromRequest(r)
	log.Debugln("SiteIdHandler request ", url)

	switch r.Method {
	case "GET":
		if site, ex := sites[url]; ex {
			writeJSON(w, site)
		} else {
			JSONError(w, fmt.Errorf("site not found %s", url))
		}

	case "PUT", "POST":

		log.Infof("POST /site/%s ", url)
		site := NewSite(url)

		site.crawlable = true
		site.crawlready = true
		acl.AddHost(url)

		// Now create the page for the new URL
		page := GetPage(url)
		page.CrawlReady = true

		log.Infoln("sending url to URLq")

		urlQ.Send(url)
		writeJSON(w, map[string]string{"saved": url})

	case "DELETE":
		h := GetHostname(url)
		DeleteSite(h)
		writeJSON(w, struct {
			Msg string
		}{
			Msg: "Deleted " + url,
		})

	default:
		JSONError(w, fmt.Errorf("unspported method "+r.Method))
	}

	// SaveSites in case of PUT / DELETE (wont hurt for GET??).
	// XXX Clearly broken, need to queue & rate limit etc..
	SaveSites()

	return
}
