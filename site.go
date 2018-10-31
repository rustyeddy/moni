package moni

import (
	"time"

	log "github.com/sirupsen/logrus"
)

// Site is basically a website wich includes API interfaces
type Site struct {
	ID      int64  `bson:"_id" json:"id"`
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
	sites Sitemap
)

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

	err := st.Put("config.json", sites)
	IfErrorFatal(err)
}

func DeleteSite(url string) {
	panic("todo")
}

// AddNewSite will create a New Site from the url, including
// verify and sanitize the url and so on.
func AddUrl(url string) {

	panic("need to do this")

	// Schedule a new crawl
	// Store the site
	log.Infof("Added new site %s ~ calling StoreSites()", url)

	Crawler.UrlQ <- url
}

// ====================================================================

// ScheduleCrawl
func (s *Site) ScheduleCrawl() {
	timer := time.AfterFunc(time.Minute*5, func() {
		s.crawlState = CrawlReady
	})
	defer timer.Stop()
}
