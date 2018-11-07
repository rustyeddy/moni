package moni

import (
	"time"

	log "github.com/sirupsen/logrus"
)

// Site is basically a website wich includes API interfaces
type Site struct {
	URL     string `bson:"url" json:"url"`
	IP      string `bson:"ip" json:"ip"`
	Health  bool   `bson:"health" json:"health"`
	Pagemap `bson:"pagemap" json:"pagemap"`

	// Crawl job info
	LastCrawled time.Time `bson:"last_crawled" json:"last_crawled"`

	nextCrawl  time.Time // Ignore this in
	crawlState int
	crawlable  bool
	crawlready bool

	*time.Timer
	*log.Entry // Ignore
}

// Sitemap
// ====================================================================
type Sitemap map[string]*Site

func initSites() (sm Sitemap) {
	if sm = ReadSites(); sm == nil {
		sm = make(map[string]*Site)
	}
	return sm
}

// NewSite will create a new *Site structure for the
// given URL
func NewSite(url string) (s *Site) {
	log.Debugln("NewSite ", url)
	s = &Site{URL: url}
	app.Sitemap[url] = s
	log.Infof("Created new site %s total sites %d", url, len(app.Sitemap))
	return s
}

func ReadSites() (sites Sitemap) {
	st := UseStore(app.Storedir)
	IfNilFatal(st)

	err := st.Get("sites.json", &sites)
	IfErrorFatal(err, "reading sites.json")

	// add sites to the acl
	for u, _ := range sites {
		app.AddHost(u)
	}
	return sites
}

func SaveSites() {
	st := UseStore(app.Storedir)
	IfNilFatal(st)

	err := st.Put("sites.json", app.Sitemap)
	IfErrorFatal(err)
}

func DeleteSite(url string) {
	if _, ex := app.Sitemap[url]; ex {
		delete(app.Sitemap, url)
	}
	SaveSites()
}

// ====================================================================
// ScheduleCrawl
func (s *Site) ScheduleCrawl() {
	timer := time.AfterFunc(time.Minute*5, func() {
		s.crawlState = CrawlReady
	})
	defer timer.Stop()
}
