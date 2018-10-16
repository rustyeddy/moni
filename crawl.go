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
	CrawlDepth int = 1
)

// Crawl will visit the given URL, and depending on configuration
// options potentially walk internal links.
func Crawl(pg *Page) {

	var site *Site

	// Create the collector and go get shit! (preserve?)
	c := colly.NewCollector(
		colly.MaxDepth(2),
	)

	c.OnRequest(func(r *colly.Request) {
		ustr := r.URL.String()
		ustr, err := NormalizeURL(ustr)
		if err != nil {
			log.Errorf("normalizing url %s => %v", r.URL, err)
			return
		}
		site = SiteFromURL(ustr)
		if site == nil {
			log.Errorln("failed to get site ", r.URL)
			return
		}
		log.Debugf("Visiting site %s ", site.URL)
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
		e.Request.Visit(newpg.URL)
	})

	c.OnResponse(func(r *colly.Response) {
		link := r.Request.URL.String()
		log.Debugln("  response from", link, "status", r.StatusCode)
		pg.StatusCode = r.StatusCode
		pg.Finish = time.Now()
		Pages[link] = pg
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Infoln("error:", r.StatusCode, err)
		pg.StatusCode = r.StatusCode
		pg.Finish = time.Now()
		Pages[r.Request.URL.String()] = pg
	})

	pg.Start = time.Now()
	c.Visit(pg.URL)
}

// CrawlOrNot will determine if the provided url is allowed to be crawled,
// and if enough time has passed before the url can be scanned again
func CrawlOrNot(urlstr string) (pi *Page) {

	allowed := ACL.IsAllowed(urlstr)
	if !allowed {
		log.Debugf("  not allowed %s add reason ..", urlstr)
		return nil
	}

	if pi = GetPage(urlstr); pi == nil {
		log.Errorf("page not found url %s", urlstr)
		return nil
	}

	if !pi.crawl {
		log.Debugf("  %s not ready to crawl ~ crawl bit off ", urlstr)
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
	ACL.AllowHost(ustr)

	// If the url has too recently been scanned we will return
	// null for the job, however a copy of the scan will is
	// available and will be returned to the caller.
	page := CrawlOrNot(ustr)
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
