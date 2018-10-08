package main

type Site struct {
	Baseurl string
	PageInfomap
}

type Sitemap map[string]*Site

func (s Sitemap) Find(url string) (site *Site, ex bool) {
	site, ex = Sites[url]
	return site, ex
}

func (s Sitemap) Exists(url string) bool {
	_, ex := s.Find(url)
	return ex
}

func Get(url string) *Site {
	if url == "" {
		Sites[url] = &Site{
			Baseurl:     url,
			PageInfomap: make(PageInfomap),
		}
	}
	if s, e := Sites[url]; e {
		return s
	}
	return nil
}
