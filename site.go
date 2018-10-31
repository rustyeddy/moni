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
	*time.Timer
	*log.Entry // Ignore
}

// Sitemap
// ====================================================================
type Sitemap map[string]*Site

var (
	sites Sitemap = make(map[string]*Site)
)

// NewSite will create a new *Site structure for the
// given URL
func NewSite(url string) (s *Site) {
	s = &Site{URL: url}
	sites[url] = s
	log.Infof("Created new site %s total sites %d", url, len(sites))
	return s
}

func ReadSites() (sites []string) {
	st := UseStore(config.Storedir)
	IfNilFatal(st)

	err := st.Get("sites.json", sites)
	IfErrorFatal(err, "reading sites.json")

	return sites
}

func SaveSites() {
	st := UseStore(config.Storedir)
	IfNilFatal(st)

	err := st.Put("sites.json", sites)
	IfErrorFatal(err)
}

func DeleteSite(url string) {
	panic("todo")
}

// ====================================================================
// ScheduleCrawl
func (s *Site) ScheduleCrawl() {
	timer := time.AfterFunc(time.Minute*5, func() {
		s.crawlState = CrawlReady
	})
	defer timer.Stop()
}
