package main

import (
	"net/url"
)

type Site struct {
	*url.URL
	Home *Page
	Pagemap

	crawl     bool
	crawlJobs []CrawlJob // List of crawl jobs
}

type Sitemap map[string]*Site

// SiteFromURL will create and index the site based on
// hostname:port.
func SiteFromURL(u *url.URL) (s *Site) {
	u2 := &url.URL{}
	u2.Host = u.Host // Can be host:port
	u2.Scheme = u.Scheme
	if u2.Scheme == "" {
		u2.Scheme = "http"
	}
	s = &Site{
		URL:     u2,
		Pagemap: make(Pagemap),
	}
	return s
}

func (s Sitemap) Find(url string) (site *Site, ex bool) {
	site, ex = Sites[url]
	return site, ex
}

func (s Sitemap) Exists(url string) bool {
	_, ex := s.Find(url)
	return ex
}

/*
func (s Sitemap) Get(urlstr string) *Site {
	if urlstr != "" {
		Sites[urlstr] = &Site{
			URL:     url.Parse(url),
			Pagemap: make(Pagemap),
		}
	}
	if s, e := Sites[url]; e {
		return s
	}
	return nil
}
*/
