package main

type Site struct {
	URL  string
	Home *Page
	Pagemap

	crawl     bool
	crawlJobs []CrawlJob // List of crawl jobs
}

type Sitemap map[string]*Site

// SiteFromURL will create and index the site based on
// hostname:port.
func SiteFromURL(ustr string) (s *Site) {
	s = &Site{
		URL:     ustr,
		Pagemap: make(Pagemap),
		crawl:   true, // assume
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
