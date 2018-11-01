package moni

import (
	"net/url"
	"strings"
	"time"

	"github.com/gocolly/colly"
	log "github.com/sirupsen/logrus"
)

// CrawlJob is created periodically to manage a crawl request
type CrawlJob struct {
	crawlId int64
	url     string  // site url for this job
	pages   []*Page // pages for this crawl
	depth   int     // crawl depth (see conni)
	time.Time
}

// Crawl will visit the given URL, and depending on configuration
// options potentially walk internal links.
//
// Order of the callbacks http://go-colly.org/docs/introduction/start/
func Crawl(pg *Page) {

	log.Infof("Crawl ~ with %s ", pg.URL)

	// create the collector and go get shit! (preserve?)
	c := colly.NewCollector(
		colly.MaxDepth(4),
		colly.DisallowedDomains("namecheap.com", "www.namecheap.com", "wordpress.org", "www.wordpress.org", "developer.wordpress.org"),
		//colly.Async(true),
	)
	// Limit parallelism to 2
	c.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 2})

	c.OnRequest(func(r *colly.Request) {
		ustr := r.URL.String()
		log.Infoln("OnRequest ", ustr)
		pg.CrawlReady = false
	})

	// OnHTML will be called when we encounter a page reference
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		if link := e.Request.AbsoluteURL(e.Attr("href")); link != "" {
			log.Infoln("OnHTML link ", link)
			// Just send the link to the URL Q for processing
			if page := GetPage(link); page != nil {
				if !page.CrawlReady {
					log.Infoln("\t skip not crawl ready")
					return
				}
			}
			urlQ.Send(link)
		}
	})

	c.OnResponse(func(r *colly.Response) {
		link := r.Request.URL.String()
		log.Infoln("  response from", link, "status", r.StatusCode)
		pg.StatusCode = r.StatusCode
		pg.Finish = time.Now()
		pg.CrawlReady = false
		pages[link] = pg
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Infoln("error:", r.StatusCode, err)
		pg.Err = err
		pg.StatusCode = r.StatusCode
		pg.Finish = time.Now()
		pg.CrawlReady = false
		link := r.Request.URL.String()
		pages[link] = pg
	})

	pg.Start = time.Now()

	log.Infoln("Visiting ", pg.URL)
	c.Visit(pg.URL)

	log.Infoln("Crawl Finished ", pg.URL)
}

// CrawlOrNot will determine if the provided url is allowed to be crawled,
// and if enough time has passed before the url can be scanned again
func CrawlOrNot(urlstr string) (pi *Page) {
	log.Infoln("crawl or not ", urlstr)
	if !acl.IsAllowed(urlstr) {
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

func NameFromURL(urlstr string) (name string) {
	u, err := url.Parse(urlstr)
	if err != nil {
		log.Errorln("NameFromURL ", err)
		return
	}

	name = u.Hostname()
	name = "crawl-" + TimeStamp() + "-" + strings.Replace(name, ".", "-", -1)
	return name
}

// FindCrawls will match a given pattern against keys in the store returning
// a list of matching crawls if there are any
func FindCrawls(pattern string) (crawls []string) {
	st := GetStore()
	crawls = st.Glob("crawl-*.json")
	return crawls
}

// GetCrawls
func GetCrawls() (crawls []string) {
	st := GetStore()
	st.Get("crawls", crawls)
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
