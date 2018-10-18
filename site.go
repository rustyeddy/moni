package main

import (
	"net/http"
)

type Site struct {
	URL    string
	IPAddr string
	Health bool
	Pagemap
}

type Sitemap map[string]*Site

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
	site, ex = Sites[url]
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
		delete(Sites, url)
	}
}

func SiteListHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, Sites)
}

func SiteIdHandler(w http.ResponseWriter, r *http.Request) {
	url := urlFromRequest(r)

	switch r.Method {
	case "GET":
		site := Sites.Get(url)
		if site != nil {
			writeJSON(w, site)
		} else {
			JSONError(w, ErrorNotFound(url+" site not found "))
		}

	case "PUT", "POST":
		Sites[url] = SiteFromURL(url)

	case "DELETE":
		Sites.Delete(url)

	default:
		JSONError(w, ErrorNotSupported("unspported method "+r.Method))
	}
}
