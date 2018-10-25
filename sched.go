package moni

import (
	"fmt"

	log "github.com/sirupsen/logrus"
)

type Scheduler struct {
	URLQ   chan string
	CrawlQ chan *Page
	ErrorQ chan error
}

var (
	sched *Scheduler
)

func GetScheduler() *Scheduler {
	if sched == nil {
		sched = &Scheduler{
			URLQ:   make(chan string, 2),
			CrawlQ: make(chan *Page, 2),
			ErrorQ: make(chan error, 2),
		}
	}
	return sched
}

func (s *Scheduler) WatchChannels() {
	for {
		select {
		case url := <-s.URLQ:

			log.Infoln("URL-Q incoming ", url)
			if pg := processURL(url, s.CrawlQ, s.ErrorQ); pg == nil {
				log.Debugln("URL-Q ~ ignoring page ", url)
			} else {
				log.Debugln("page -> crawlQ ", pg.URL)
				s.CrawlQ <- pg
				log.Debugln("        crawlQ finished ")
			}

		case npg := <-s.CrawlQ:
			Crawl(npg)
			storePageCrawl(npg)
		}
	}
}

func processURL(ustr string, crawlq chan *Page, errch chan error) *Page {

	log.Infoln("processURL ", ustr)

	// Normalize the URL and fill in a scheme it does not exist
	ustr, err := NormalizeURL(ustr)
	if err != nil {
		log.Infoln("processURL Normalize failed ", ustr)
		errch <- fmt.Errorf("I had a problem with the url %v", ustr)
		return nil
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
		log.Infoln("processURL rejected failed ", ustr)
		errch <- fmt.Errorf("url rejected %s", ustr)
		return nil
	}

	log.Infoln("returning from process URL with page ", page.URL)
	return page
}
