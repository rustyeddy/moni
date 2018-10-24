package moni

import (
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// Site is basically a website wich includes API interfaces
type Site struct {
	URL    string
	IP     string
	Health bool
	Pagemap
}

type Sitemap map[string]*Site

var (
	sites Sitemap = make(Sitemap, 10)
)

func GetSites() Sitemap {
	if sites == nil || len(sites) < 1 {
		st := GetStorage()
		if _, err := st.FetchObject("sites", sites); err != nil {
			log.Errorf(" failed to read saved 'sites' %v", err)
			sites = make(Sitemap)
		} else {
			log.Infof("retrived sites from the filesystem %+v", sites)
		}
	}
	return sites
}

// StoreSites will attempt to store our memory version of
// the Sitemap to a file.  Hope it all works out, we will get
// a log message if there is a problem
func StoreSites() {
	if st := GetStorage(); st != nil {
		if obj, err := st.StoreObject("sites", sites); err != nil {
			log.Errorf("failed StoreObject sites %v", err)
		} else {
			log.Infof("sites stored object: %+v\n", obj)
		}
	}
}

// AddNewSite will create a New Site from the url, including
// verify and sanitize the url and so on.
func AddNewSite(url string) *Site {
	sites := GetSites()
	s := &Site{
		URL:     url,
		Pagemap: make(Pagemap),
	}
	sites[url] = s

	log.Infof("Added new site %s ~ calling StoreSites()", url)
	StoreSites()
	return s
}

func (s Sitemap) Find(url string) (site *Site, ex bool) {
	site, ex = s[url]
	return site, ex
}

func (s Sitemap) Get(url string) (site *Site) {
	if site, ex := s.Find(url); ex {
		return site
	}
	return nil
}

func (s Sitemap) Exists(url string) bool {
	_, ex := s.Find(url)
	return ex
}

func (s Sitemap) Delete(url string) {
	if _, ex := s[url]; ex {
		delete(s, url)
	}
}

func (s Sitemap) Store() {
	st := GetStorage()
	if _, err := st.StoreObject("sites", s); err != nil {
		log.Errorf("failed saving sites %v", err)
	}
}

func (s *Sitemap) Fetch() {
	st := GetStorage()
	if _, err := st.FetchObject("sites", s); err != nil {
		log.Errorf("failed to fetch sites %v", err)
	}
}

// REST API
// ====================================================================

func SiteListHandler(w http.ResponseWriter, r *http.Request) {
	sites := GetSites()
	if sites == nil || len(sites) < 1 {
		fmt.Fprintf(w, "no sites ")
	}
	writeJSON(w, sites)
}

func SiteIdHandler(w http.ResponseWriter, r *http.Request) {
	url := urlFromRequest(r)
	sites := GetSites()

	switch r.Method {
	case "GET":
		site := sites.Get(url)
		if site != nil {
			writeJSON(w, site)
		} else {
			JSONError(w, ErrorNotFound(url+" site not found "))
		}

	case "PUT", "POST":
		AddNewSite(url)

	case "DELETE":
		sites.Delete(url)

	default:
		JSONError(w, ErrorNotSupported("unspported method "+r.Method))
	}
	writeJSON(w, "OK")
}

// WebUI
// ====================================================================
type SitesCard struct {
	*Card
	*Sitemap
}

func GetSitesCard() (sc *SitesCard) {
	var s Sitemap
	if s = GetSites(); s == nil {
		s = make(Sitemap)
	}
	c := &Card{}
	sc = &SitesCard{
		Card:    c,
		Sitemap: &s,
	}
	return sc
}
