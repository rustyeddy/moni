package main

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
	Visited    Pagemap = make(Pagemap)
	CrawlDepth int     = 1
)

// Crawl will visit the given URL, and depending on configuration
// options potentially walk internal links.
func Crawl(pg *Page) {

	// Create the collector and go get shit! (preserve?)
	c := colly.NewCollector(
		colly.MaxDepth(2),
	)

	c.OnRequest(func(r *colly.Request) {
		ustr := r.URL.String()
		u, err := NormalizeURL(ustr)
		if err != nil {
			log.Errorf("normalizing url %s => %v", r.URL, err)
			return
		}
		site := SiteFromURL(u)
		if site == nil {
			log.Errorln("failed to get site ", r.URL)
			return
		}
		log.Infof("Visiting site %s ", site.URL)
	})

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Request.AbsoluteURL(e.Attr("href"))
		if link == "" {
			return
		}

		u, err := NormalizeURL(link)
		if err != nil {
			log.Errorf("expected url got error %v", err)
			return
		}
		hp := Hostport(u)

		newpg := CrawlOrNot(u)
		if newpg == nil {
			// The link has been filtered for one
			// reason or another, we will move along
			log.Debugf("  ignoring %s ", hp)
			pg.Ignored[hp]++
			return
		}
		pg.Links[hp] = newpg
		e.Request.Visit(newpg.URL)
	})

	c.OnResponse(func(r *colly.Response) {
		link := r.Request.URL
		log.Infoln("  response from", link, "status", r.StatusCode)
		pg.StatusCode = r.StatusCode
		pg.Finish = time.Now()
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Infoln("error:", r.StatusCode, err)
		pg.StatusCode = r.StatusCode
		pg.Finish = time.Now()
	})

	pg.Start = time.Now()
	c.Visit(pg.URL)
}

// CrawlOrNot will determine if the provided url is allowed to be crawled,
// and if enough time has passed before the url can be scanned again
func CrawlOrNot(url *url.URL) (pi *Page) {

	// Check if we will allow crawling this hostname
	hp := Hostport(url)

	allowed := ACL.IsAllowed(hp)
	if !allowed {
		log.Debugf("  not allowed %s add reason ..", hp)
		return nil
	}

	if pi = GetPage(url); pi == nil {
		log.Errorf("page not found url %s hostport %s", url.String(), hp)
		return nil
	}

	if !pi.crawl {
		log.Debugf("  %s not ready to crawl ~ crawl bit off ", hp)
		return nil
	}
	return pi
}

// CrawlHandler will handle incoming HTTP request to crawl a URL
func CrawlHandler(w http.ResponseWriter, r *http.Request) {

	// Prepare for Execution
	// ====================================================================

	// Get host URL from mux.vars
	vars := mux.Vars(r)

	// Normalize the URL and fill in a scheme it does not exist
	url, err := NormalizeURL(vars["url"])
	if err != nil {
		fmt.Fprintf(w, "I had a problem with the url %v", url)
		return
	}

	// Add this host to the allowed host name, that way CrawlOrNot
	// will not reject the hostname.  TODO: set a config to allow
	// this to be turned on or off.
	hname := Hostport(url)
	ACL.AllowHost(hname)

	// CrawlOrNot will determine if we are permited to crawl this
	// link. If we are permitted, is the link ready to be crawled.
	// CrawlOrNot will figure that our and return a job ready to
	// be scheduled (or just played).
	//
	// If the url has too recently been scanned we will return
	// null for the job, however a copy of the scan will is
	// available and will be returned to the caller.
	ustr := url.String()
	page := CrawlOrNot(url)
	if page == nil {
		log.Errorf("url rejected %s", ustr)
		fmt.Fprintf(w, "url rejected %s", ustr)
		return
	}

	// ^^^ This is where the job gets queued (written to the job channel) ^^^
	// vvv This is where the (next free?) Crawler grabs a crawl job
	Crawl(page)

	// Cache the results ...  We'll replace any '/'
	// in the URL with '-' and store the results in
	// a cache.
	name := nameFromURL(ustr)
	obj, err := Storage.StoreObject(name, page)
	if err != nil {
		log.Errorln("Failed to create local store")
		fmt.Fprintf(w, "Internal Error %s", ustr)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(obj.Buffer)
}

func nameFromURL(urlstr string) (name string) {
	u, err := url.Parse(urlstr)
	if err != nil {
		log.Errorln(err)
		return
	}
	name = u.Hostname()
	name = strings.Replace(name, ".", "-", -1)
	return name
}
