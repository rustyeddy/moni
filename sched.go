package moni

import (
	"fmt"
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
			URLQ:   make(chan string),
			CrawlQ: make(chan *Page),
			ErrorQ: make(chan error),
		}
	}
	return sched
}

func (s *Scheduler) WatchChannels() {
	for {
		select {
		case url := <-s.URLQ:
			s.CrawlQ <- processURL(url, s.CrawlQ, s.ErrorQ)
		case npg := <-s.CrawlQ:
			Crawl(npg)
			storePageCrawl(npg)
		}
	}
}

func processURL(ustr string, crawlq chan *Page, errch chan error) *Page {

	// Normalize the URL and fill in a scheme it does not exist
	ustr, err := NormalizeURL(ustr)
	if err != nil {
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
		errch <- fmt.Errorf("url rejected %s", ustr)
	}
	return page
}
