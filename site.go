package moni

import (
	"encoding/json"
	"io/ioutil"
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
	return s
}

func ReadSites() (sites Sitemap) {
	var (
		buf []byte
		err error
	)
	st := GetStore()
	IfNilFatal(st)

	path := st.PathFromName("sites.json")
	if buf, err = ioutil.ReadFile(path); err != nil {
		log.Errorf("read index %s failed %v", path, err)
		return nil
	}

	sites = make(Sitemap)
	switch st.ContentType {
	case "application/json":
		if err = json.Unmarshal(buf, &sites); err != nil {
			IfErrorFatal(err, "get failed marshaling json "+path)
		}
	default:
		panic("did not expect this")
	}
	return sites
}

func SaveSites() {
	st := GetStore()
	IfNilFatal(st)

	err := st.Put("sites.json", sites)
	IfErrorFatal(err)
}

func DeleteSite(url string) {
	if _, ex := sites[url]; ex {
		delete(sites, url)
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
