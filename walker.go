package main

import (
	"encoding/json"
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
	var urls []string

	// make the page more convinient
	p := w.Page

	nilPanic(w)
	c := colly.NewCollector()

	log.Infof("Visiting page %s", p.URL.String())

	// Setup all the collbacks
	c.OnHTML("a", func(e *colly.HTMLElement) {
		refurl := e.Attr("href")
		link := e.Request.AbsoluteURL(refurl)
		p.Links[link] = append(p.Links[link], e.Text)
	})

	c.OnRequest(func(r *colly.Request) {
		pages[*r.URL] = NewPage(*r.URL)
		log.Infoln("Request ", r.URL)
	})

	c.OnResponse(func(r *colly.Response) {
		p.RespTime = time.Now()
		log.Infoln("Response ", r.Request.URL)
	})

	c.OnScraped(func(r *colly.Response) {

		// The page scrape has completed
		log.Debugf("\tLinks: %s", p.URL.String())
		for ustr, _ := range p.Links {
			log.Debugf("\t~> %s", ustr)
			if config.Recurse {
				urls = append(urls, ustr)
			}
		}
	})

	// Start the walk
	p.ReqTime = time.Now()
	c.Visit(p.String())

	// Walk is complete
	p.RespTime = time.Now()
	p.Elapsed = p.RespTime.Sub(p.ReqTime)

	links := []string{}
	for l, _ := range p.Links {
		links = append(links, l)
	}

	resp := struct {
		Status     string        `json:"status"`
		ReturnCode int           `json:"returnCode"`
		URL        string        `json:"url"`
		Links      []string      `json:"links"`
		Elapsed    time.Duration `json:"elapsed"`
	}{
		Status:     "",
		ReturnCode: 200,
		URL:        p.URL.String(),
		Links:      links,
		Elapsed:    p.Elapsed,
	}

	//fmt.Fprintf(w, "Responding From WALK %+v\n", resp)
	json.NewEncoder(w.w).Encode(resp)
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
				Page: page,
			}
			walker.URL = u
			walker.Walk()
		}
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
