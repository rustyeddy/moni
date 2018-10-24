package moni

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

var (
	CrawlDepth int = 1
)

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
			log.Errorf("normalizing url %s => %v", r.URL, err)
			return
		}
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

		if newpg.crawl {
			e.Request.Visit(newpg.URL)
			newpg.crawl = false
		}
	})

	c.OnResponse(func(r *colly.Response) {
		link := r.Request.URL.String()
		log.Debugln("  response from", link, "status", r.StatusCode)
		pg.StatusCode = r.StatusCode
		pg.Finish = time.Now()
		pg.crawl = false
		pages[link] = pg
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Infoln("error:", r.StatusCode, err)
		pg.Err = err
		pg.StatusCode = r.StatusCode
		pg.Finish = time.Now()
		pg.crawl = false
		link := r.Request.URL.String()
		pages[link] = pg
	})

	pg.Start = time.Now()
	c.Visit(pg.URL)
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

	if !pi.crawl {
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

// ServiceHandlers
// ========================================================================

// CrawlHandler will handle incoming HTTP request to crawl a URL
func CrawlHandler(w http.ResponseWriter, r *http.Request) {

	// Prepare for Execution
	// Extract the url(s) that we are going to walk
	vars := mux.Vars(r)
	ustr := vars["url"]

	// Normalize the URL and fill in a scheme it does not exist
	ustr, err := NormalizeURL(ustr)
	if err != nil {
		fmt.Fprintf(w, "I had a problem with the url %v", ustr)
		return
	}

	// This conversion back to string is necessary and simple domain
	// name like "example.com" will be placeded in the url.URL.Path
	// field instead of the Host field.  However url.String() makes
	// everything right.
	accessList.AllowHost(ustr)

	// If the url has too recently been scanned we will return
	// null for the job, however a copy of the scan will is
	// available and will be returned to the caller.
	page := CrawlOrNot(ustr)
	if page == nil {
		log.Errorf("url rejected %s", ustr)
		fmt.Fprintf(w, "url rejected %s", ustr)
		return
	}

	// ~~~ This is where the (next free?) Crawler grabs a crawl job ~~~
	Crawl(page)

	// Write the results back to the caller
	writeJSON(w, page)

	// Cache the results ...  We'll replace '/' with '-' and
	// store the results in the cache store.
	storePageCrawl(page)
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

// CrawlListHandler will return a list of all recent crawls.
// As stored in our storage (json) file.
func CrawlListHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, GetCrawls())
}

// CrawlIdHandler will return a list of all recent crawls.
// As stored in our storage (json) file.
func CrawlIdHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the url(s) that we are going to walk
	vars := mux.Vars(r)
	cid := vars["cid"]
	st := GetStorage()

	page := new(Page)
	_, err := st.FetchObject(cid, page)
	if err != nil {
		JSONError(w, err)
		return
	}
	writeJSON(w, page)
}
