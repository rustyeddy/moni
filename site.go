package moni

import (
	"time"
)

// Site is basically a website wich includes API interfaces
type Site struct {
	URL    string
	IP     string
	Health bool
	Pagemap

	// Crawl job info
	lastCrawled time.Time
	nextCrawl   time.Time
	CrawlState  int

	*time.Timer
}

type Sitemap map[string]*Site

var (
	sites Sitemap = make(Sitemap)
)

func (app *App) FetchSites() Sitemap {
	if _, err := storage.FetchObject("sites", sites); err != nil {
		log.Errorf(" failed to read saved 'sites' %v", err)
		sites = make(Sitemap)
	} else {
		log.Infof("retrived Sites from the filesystem %+v", sites)
	}
	return sites
}

// StoreSites will attempt to store our memory version of
// the Sitemap to a file.  Hope it all works out, we will get
// a log message if there is a problem
func (app *App) StoreSites() {
	if obj, err := storage.StoreObject("sites", sites); err != nil {
		log.Errorf("failed StoreObject Sites %v", err)
	} else {
		log.Infof("Sites stored object: %+v\n", obj)
	}
}

// AddNewSite will create a New Site from the url, including
// verify and sanitize the url and so on.
func AddNewSite(url string) *Site {
	s := &Site{
		URL:     url,
		Pagemap: make(Pagemap),
	}
	sites[url] = s

	// Schedule a new crawl
	// Store the site
	log.Infof("Added new site %s ~ calling StoreSites()", url)

	Crawler.UrlQ <- url

	// This should not cause any problems, que no?
	go app.StoreSites()
	return s
}

// RemoveSite represented by the URL from the list of sites to manage
func RemoveSite(url string) {

	// Unschedule the site from the crawler
	log.Infoln("Deleting URL ", url)
	sites.Delete(url)
}

// ====================================================================

// ScheduleCrawl
func (s *Site) ScheduleCrawl() {
	timer := time.AfterFunc(time.Minute*5, func() {
		s.CrawlState = CrawlReady
	})
	defer timer.Stop()
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
