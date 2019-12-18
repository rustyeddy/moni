package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/gocolly/colly"
	log "github.com/sirupsen/logrus"
)

// Walk response manage the walk of this particular host
type Walker struct {
	*url.URL `json:"url"`
	*Page    `json:"page"`
	Status   int
	w        http.ResponseWriter `json:"-"`
}

// Walk the given page, setting the links and responding to request
func (w *Walker) Walk() {
	c := colly.NewCollector()

	log.Infof("Visiting page %+v", w)

	// Setup all the collbacks
	c.OnHTML("a", func(e *colly.HTMLElement) {
		refurl := e.Attr("href")
		link := e.Request.AbsoluteURL(refurl)
		w.Page.Links[link] = append(w.Page.Links[link], e.Text)
	})

	c.OnRequest(func(r *colly.Request) {
		pages[*r.URL] = NewPage(*r.URL)
		log.Infoln("Request ", r.URL)
	})

	c.OnResponse(func(r *colly.Response) {
		w.TimeStamp = NewTimestamp()
		log.Infoln("Response ", r.Request.URL)
	})

	c.OnScraped(func(r *colly.Response) {
		// The page scrape has completed, set the response time
		w.Page.SetResponseTime(time.Now())
		w.Page.TimeStamps = append(w.Page.TimeStamps, w.Page.TimeStamp)

		// Now print some interactive user friendly stuff
		log.Infof("    Links: %s", w.Page.URL.String())
		for ustr, _ := range w.Page.Links {
			log.Infof("\t~> %s", ustr)
		}
	})

	// Start the walk
	w.Page.TimeStamp = NewTimestamp()
	c.Visit(w.Page.String())

	// Walk is complete, SetResponse time will set the time as
	// advertised, but it will also calculate the elapsed time
	// as a side effect.
	links := []string{}
	for l, _ := range w.Page.Links {
		links = append(links, l)
	}

	if w.w != nil {
		json.NewEncoder(w.w).Encode(w.Page)
	} else if config.Daemon {
		log.Infof("%+v", w.Page)
	} else {
		fmt.Println(w.Page.PageString())
	}
}

func processURLs(urls []string, w http.ResponseWriter) {
	for _, ustr := range urls {
		processURL(ustr, w)
	}
}

func processURL(ustr string, w http.ResponseWriter) {
	log.Infof("Walking %s\n", ustr)

	if u := scrubURL(ustr); u != nil {
		if page := GetPage(*u); page != nil {
			log.Infof("got page: %+v - let's walk...\n", page)
			walker := Walker{
				w:    w,
				URL:  u,
				Page: page,
			}
			walker.Walk()
		} else {
			log.Warnf("\tprocessURL: page rejected %v", *u)
		}
	} else {
		log.Warnf("\tprocessURL: URL failed %s\n", ustr)
	}
}

func scrubURL(urlstr string) (u *url.URL) {
	var err error

	log.Infoln("scrubURL with ", urlstr)

	u, err = url.Parse(urlstr)
	errPanic(err)

	if u.Scheme == "" {
		u.Scheme = "http"
	}

	// Recon struct
	if u, err = url.Parse(u.String()); err != nil {
		return nil
	}

	// if this hostname exists in the acl set as false,
	// we will just return
	if f, ex := acl[u.Host]; ex && f == false {
		return nil
	}
	return u
}
