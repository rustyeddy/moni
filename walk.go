package main

import (
	"net/url"
	"sync"
	"time"

	"github.com/gocolly/colly"
	log "github.com/sirupsen/logrus"
)

// Walk the given page, setting the links and responding to request
func (p *Page) Walk() {
	c := colly.NewCollector()

	log.Infof("Walking page %s", p.URL.String())

	// Setup all the collbacks
	c.OnHTML("a", func(e *colly.HTMLElement) {
		refurl := e.Attr("href")
		link := e.Request.AbsoluteURL(refurl)
		p.Links[link]++ //append(p.Links[link], e.Text)

		if pg := processURL(link); pg != nil {
			c.Visit(pg.URL.String())
		}
	})

	c.OnRequest(func(r *colly.Request) {
		log.Infoln("Request ", r.URL)
	})

	c.OnResponse(func(r *colly.Response) {
		p.TimeStamp = NewTimestamp()
		log.Infoln("Response ", r.Request.URL)
	})

	c.OnScraped(func(r *colly.Response) {
		// The page scrape has completed, set the response time
		p.SetResponseTime(time.Now())
		p.TimeStamps = append(p.TimeStamps, p.TimeStamp)

		// Now print some interactive user friendly stuff
		log.Infof("    Links: %s", p.URL.String())
		for ustr, _ := range p.Links {
			log.Infof("\t~> %s", ustr)
		}
	})

	// Start the walk
	p.TimeStamp = NewTimestamp()
	c.Visit(p.String())
	log.Infof("    Elaspsed time: %v", p.Elapsed)
}

func doWatcher(wQ chan *Page, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case p := <-wQ:
			log.Infof("  wQ channel recieved page: %+v", p)
			p.Walk()
		}
	}
}

func processURL(urlstr string) (pg *Page) {
	if u := scrubURL(urlstr); u != nil {
		if site := GetSite(u); site != nil {
			if pg := site.GetPage(*u); pg != nil {
				return pg
			}
		}
	}
	return nil
}

func scrubURL(urlstr string) (u *url.URL) {
	var err error

	log.Infoln("scrubURL with ", urlstr)

	u, err = url.Parse(urlstr)
	errPanic(err)

	if u.Scheme == "" {
		u.Scheme = "http"

		u, err = url.Parse(u.String())
		if err != nil {
			log.Errorf("Failed to reconstruct URL %+v, %v", u, err)
			return nil
		}
	}

	// if this hostname exists in the acl set as false,
	// we will just return
	f, ex := acl[u.Host]
	if ex == false || f == false {
		return nil
	}
	return u
}
