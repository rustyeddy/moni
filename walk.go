package main

import (
	"net/url"
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
		lstr := e.Request.AbsoluteURL(refurl)
		if lstr != "" {
			p.Links[lstr] = NewLink(lstr)
		}
		walkQ <- lstr
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
			log.Infof("\t~ %s", ustr)
		}
	})

	// Start the walk
	p.TimeStamp = NewTimestamp()

	log.Infof("Starting visit for %s", p.String())
	c.Visit(p.String())

	// Schedule new visit for website
	//scheduleVisit(p)
	walkQ <- p.URL.String()

	log.Infof("    Elaspsed time: %v", p.Elapsed)
}

func scheduleVisit(p *Page) {
	
	// if we have a walktimer a visit has already been scheduled
	if p.WalkTimer == nil {
		log.Debugf("visit already scheduled, ignore: %s", p.URL.String())
		p.WalkTimer = time.NewTimer(time.Minute * time.Duration(config.Wait))
		p.Walk()

		go func() {
			for {
				<-p.WalkTimer.C
				p.Walk()
			}
		}()
	}
}

func processURL(urlstr string) (pg *Page) {
	if u := scrubURL(urlstr); u != nil {
		if site := GetSite(u.String()); site != nil {
			if pg := site.GetPage(*u); pg != nil {
				return pg
			}
		}
	}
	return nil
}

func scrubURL(urlstr string) (u *url.URL) {
	var err error

	u, err = url.Parse(urlstr)
	if err != nil {
		log.Infof("\turlstring is bad %s...", urlstr)
		return nil
	}

	if u.Scheme == "" {
		u.Scheme = "http"

		u, err = url.Parse(u.String())
		if err != nil {
			log.Infof("Failed to reconstruct URL %+v, %v", u, err)
			return nil
		}
	}

	// if this hostname exists in the acl set as false,
	// we will just return
	ok := acl.Allow(u.Host)
	if !ok {
		log.Infof("\trejecting url do to acl %s", u.Host)
		return nil
	}
	return u
}
