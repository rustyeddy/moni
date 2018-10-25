package moni

import (
	"net/url"
	"strings"
	"time"

	"github.com/gocolly/colly"
	log "github.com/sirupsen/logrus"
)

var (
	CrawlDepth int = 1
)

func crawlWatcher(chpg chan *Page, errch chan error) {
	for {
		pg := <-chpg
		log.Infoln("CrawlWatcher")
		Crawl(pg)

		st := GetStorage()
		st.StoreObject(pg.URL, pg)
	}
}

// Crawl will visit the given URL, and depending on configuration
// options potentially walk internal links.
func Crawl(pg *Page) {

	// Create the collector and go get shit! (preserve?)
	c := colly.NewCollector(
		colly.MaxDepth(4),
	)

	c.OnRequest(func(r *colly.Request) {
		ustr := r.URL.String()
		ustr, err := NormalizeURL(ustr)
		if err != nil {
			pg.CrawlState = CrawlErrored
			log.Errorf("normalizing url %s => %v", r.URL, err)
			return
		}
		pg.CrawlState = CrawlRequestSent
	})

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Request.AbsoluteURL(e.Attr("href"))
		if link == "" {
			return
		}

		// Parsing and converting seems necessary
		ustr, err := NormalizeURL(link)
		if err != nil {
			log.Errorf("expected url got error %v", err)
			return
		}

		newpg := CrawlOrNot(ustr)
		if newpg == nil {
			// The link has been filtered for one
			// reason or another, we will move along
			log.Debugf("  ignoring %s ", ustr)
			pg.Ignored[ustr]++
			return
		}
		pg.Links[ustr] = newpg

		if newpg.CrawlState == CrawlReady {
			e.Request.Visit(newpg.URL)

		}
	})

	c.OnResponse(func(r *colly.Response) {
		link := r.Request.URL.String()
		log.Debugln("  response from", link, "status", r.StatusCode)
		pg.StatusCode = r.StatusCode
		pg.Finish = time.Now()
		pg.CrawlState = CrawlResponseRecieved
		pages[link] = pg
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Infoln("error:", r.StatusCode, err)
		pg.Err = err
		pg.StatusCode = r.StatusCode
		pg.Finish = time.Now()
		pg.CrawlState = CrawlErrored
		link := r.Request.URL.String()
		pages[link] = pg
	})

	pg.Start = time.Now()
	c.Visit(pg.URL)

	log.Infoln("Crawl Finished ", pg.URL)
}

// CrawlOrNot will determine if the provided url is allowed to be crawled,
// and if enough time has passed before the url can be scanned again
func CrawlOrNot(urlstr string) (pi *Page) {

	allowed := ACL().IsAllowed(urlstr)
	if !allowed {
		log.Debugf("  not allowed %s add reason ..", urlstr)
		return nil
	}

	if pi = PageFromURL(urlstr); pi == nil {
		log.Errorf("page not found url %s", urlstr)
		return nil
	}

	if pi.CrawlState != CrawlReady {
		log.Debugf("  %s not ready to crawl ~ crawl bit off ", urlstr)
		return nil
	}
	return pi
}

// Cache the page if there is an error we just won't have a
// cached page and will need to refetch.  Can get ugly if
// the page is fetched a lot.
func storePageCrawl(pg *Page) {
	name := NameFromURL(pg.URL)
	st := GetStorage()
	_, err := st.StoreObject(name, pg)
	if err != nil {
		log.Errorln("Failed to create local store")
	}
}

func NameFromURL(urlstr string) (name string) {
	u, err := url.Parse(urlstr)
	if err != nil {
		log.Errorln(err)
		return
	}

	name = u.Hostname()
	name = "crawl-" + TimeStamp() + "-" + strings.Replace(name, ".", "-", -1)
	return name
}

// GetCrawls
func GetCrawls() []string {
	st := GetStorage()
	pat := "crawl-"
	patlen := len(pat)

	crawls, _ := st.FilterNames(func(name string) string {
		if len(name) < patlen {
			return ""
		}
		if name[0:patlen] == pat {
			return name
		}
		return ""
	})
	if crawls == nil {
		crawls = []string{}
	}
	return crawls
}
