package moni

import (
	"time"

	log "github.com/sirupsen/logrus"
)

type Q interface {
	Send(msg string)
	Watch()
}

// ChannelBase
type QBase struct {
	Q    chan string
	qlen int

	Incoming, Outgoing int64
}

//			      URL Cleansing and Normalization
// ====================================================================

// URLChannel is used to accept incoming url to be sanitized,
// normalized and processed.
type URLQ struct {
	QBase
}

func NewURLQ() *URLQ {
	q := new(URLQ)
	q.qlen = 2
	q.Q = make(chan string, q.qlen)
	return q
}

// Watch URLChannel is meant to be a go routine consuming
// incoming Url channel requests
func (q *URLQ) Watch() {
	var ts time.Time = time.Now()
	for {
		log.Infof("Watching the URL channel at ..")

		select {
		case url := <-q.Q:
			var pg *Page

			q.Incoming++
			log.Infof("  <~~ urlQ <~~ incoming[%d] %s at %v ", q.Incoming, url, time.Since(ts))

			log.Infof("     must processURL at %v ", time.Since(ts))
			if pg = processURL(url); pg == nil {
				log.Warnf("  !!! process URL failed %s skipping... ", url)
				continue
			}

			// if pg.CrawlReady {
			// 	crawlQ.Send(pg)
			// } else {
			// 	log.Infof("  !!! page is not crawlQ ready skipping %s at %v", url, time.Since(ts))
			// }
			log.Debugf("     urlQ incoming[%d] complete: %s", q.Incoming, url)
		}
	}
}

func (q *URLQ) Send(url string) {
	ts := time.Now()
	log.Infof("  ~~> urlQ ~~> send page %s crawl channel", url)
	q.Q <- url
	log.Infof("      urlQ send complete %s, %v", url, time.Since(ts))
}

//			      Crawl Channel
// ====================================================================

type CrawlQ struct {
	Q chan *Page
	QBase
}

func NewCrawlQ() *CrawlQ {
	q := new(CrawlQ)
	q.qlen = 2
	q.Q = make(chan *Page, q.qlen)
	return q
}

func (q *CrawlQ) Watch() {
	ts := time.Now()

	for {
		log.Infof("watching the crawl channel at %v ", ts)
		select {
		case page := <-q.Q:

			ts = time.Now()
			log.Infof("  <~~ crawlQ <~~ incoming for %s at %v", page.URL, ts)

			Crawl(page)
			log.Infof("      crawlQ finished for %s after %v", page.URL, time.Since(ts))
		}
	}
}

func (q *CrawlQ) Send(pg *Page) {
	ts := time.Now()

	log.Infof("  ~~> crawlQ ~~> sending to crawlQ %s to crawl channel at %v ", pg.URL, ts)
	q.Q <- pg
	log.Infof("      crawlQ send complete for %s at %v", pg.URL, time.Since(ts))
	q.Outgoing++
}

//			      Save Channel
// ====================================================================

type SaveQ struct {
	Q chan *Page
	QBase
}

func NewSaveQ() *SaveQ {
	q := new(SaveQ)
	q.qlen = 2 // arbitrary
	q.Q = make(chan *Page, q.qlen)
	return q
}

func (q *SaveQ) Watch() {
	ts := time.Now()

	for {
		log.Infof("watching SaveQ at %v ", ts)
		select {
		case page := <-q.Q:
			ts = time.Now()
			log.Infof("  <~~ saveQ <~~ incoming request for %s at %v", page.URL, ts)
			StorePage(page)
			log.Infof("      saveQ complete for %s after %v", page.URL, time.Since(ts))
		}
	}
}

func (q *SaveQ) Send(pg *Page) {
	ts := time.Now()

	log.Infof("  ~~> saveQ ~~> send page %s to be saved %v ", pg.URL, ts)
	q.Q <- pg
	log.Infof("      saveQ send complete for %s at %v", pg.URL, time.Since(ts))
}
