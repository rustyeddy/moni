package moni

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

// Crawler represents the crawling
type CrawlDispatcher struct {
	*App
	UrlQ   chan string // public used externally to submit urls
	crawlQ chan *Page
	saveQ  chan *Page
	errQ   chan error

	qsize int
}

// CrawlJob is created periodically to manage a crawl request
type CrawlJob struct {
	crawlId int64
	url     string  // site url for this job
	pages   []*Page // pages for this crawl
	depth   int     // crawl depth (see conni)
	time.Time
}

var (
	CrawlDepth int = 2
	Crawler    *CrawlDispatcher
)

func init() {
	Crawler = NewCrawler()
}

// NewCrawler will handle scheduling all call requests
func NewCrawler() (crawler *CrawlDispatcher) {
	cr := &CrawlDispatcher{
		qsize: 2,
	}
	cr.UrlQ = make(chan string, cr.qsize)
	cr.crawlQ = make(chan *Page, cr.qsize)
	cr.saveQ = make(chan *Page, cr.qsize)
	cr.errQ = make(chan error, cr.qsize)
	Crawler = cr
	return Crawler
}

func (cr *CrawlDispatcher) WatchChannels() {
	for {
		log.Infoln("URLQ Channel Watcher waiting for URL ... ")
		ts := time.Now()

		select {
		case url := <-cr.UrlQ:
			log.Infof("urlChan recieved %s ~ %v ", url, time.Since(ts))
			if pg := processURL(url); pg != nil {
				cr.crawlQ <- pg
			} else {
				cr.errQ <- fmt.Errorf("url %s has errored ", url)
			}

		case page := <-cr.crawlQ:

			cr.Crawl(page)

		case page := <-cr.saveQ:
			storage.StoreObject(page.URL, page)

		case err := <-cr.errQ:
			log.Error(err)
		}
	}
}

// Crawl will visit the given URL, and depending on configuration
// options potentially walk internal links.
func (cr *CrawlDispatcher) Crawl(pg *Page) {

	// Create the collector and go get shit! (preserve?)
	c := colly.NewCollector(
		colly.MaxDepth(4),
	)

	c.OnRequest(func(r *colly.Request) {
		ustr := r.URL.String()

		cr.UrlQ <- ustr
	})

	// OnHTML will be called when we encounter a page reference
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		if link := e.Request.AbsoluteURL(e.Attr("href")); link != "" {
			// Just send the link to the URL Q for processing
			cr.UrlQ <- link
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
	log.Infoln("crawl or not ", urlstr)
	allowed := ACL().IsAllowed(urlstr)
	if !allowed {
		log.Infof("  not allowed %s add reason ..", urlstr)
		return nil
	}

	if pi = PageFromURL(urlstr); pi == nil {
		log.Errorf("page not found url %s", urlstr)
		return nil
	}

	if pi.CrawlState != CrawlReady {
		log.Infof("  %s not ready to crawl ~ crawl bit off ", urlstr)
		return nil
	}
	return pi
}

// Cache the page if there is an error we just won't have a
// cached page and will need to refetch.  Can get ugly if
// the page is fetched a lot.
func storePageCrawl(pg *Page) {
	name := NameFromURL(pg.URL)
	storage.StoreObject(name, pg)
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
	pat := "crawl-"
	patlen := len(pat)

	crawls, _ := storage.FilterNames(func(name string) string {
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

// GetTimeStamp returns a timestamp in a modified RFC3339
// format, basically remove all colons ':' from filename, since
// they have a specific use with Unix pathnames, hence must be
// escaped when used in a filename.
func TimeStamp() string {
	ts := time.Now().UTC().Format(time.RFC3339)
	return strings.Replace(ts, ":", "", -1) // get rid of offesnive colons
}
