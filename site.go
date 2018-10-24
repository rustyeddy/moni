package moni

import (
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type Site struct {
	URL    string
	IP     string
	Health bool
	Pagemap
}

type Sitemap map[string]*Site

type SitesCard struct {
	*Card
	*Sitemap
}

var (
	sites Sitemap = make(Sitemap, 10)
)

func GetSites() Sitemap {
	if sites == nil || len(sites) < 1 {
		st := GetStorage()
		if _, err := st.FetchObject("sites", sites); err != nil {
			log.Errorf(" failed to read saved 'sites' %v", err)
			sites = make(Sitemap)
		}
	}
	return sites
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

// SiteFromURL will create and index the site based on
// hostname:port.
func SiteFromURL(ustr string) (s *Site) {
	s = &Site{
		URL:     ustr,
		Pagemap: make(Pagemap),
	}
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
		s := SiteFromURL(url)
		fmt.Printf("sites: %+v\n", sites)
		sites[url] = s

	case "DELETE":
		sites.Delete(url)

	default:
		JSONError(w, ErrorNotSupported("unspported method "+r.Method))
	}
	writeJSON(w, "OK")
}
