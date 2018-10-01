package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gocolly/colly"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type pageInfo struct {
	StatusCode int
	Links      map[string]int
	Start      time.Time
	End        time.Time
}

func handleWalk(w http.ResponseWriter, r *http.Request) {

	// Using Gorilla mux /walk/{url}
	vars := mux.Vars(r)
	url := vars["url"]

	// TODO: check for "http://" prefix
	url = "http://" + url
	log.Infoln("visiting", url)

	// Create the collector and go get shit!
	c := colly.NewCollector(
		colly.MaxDepth(config.Depth),
	)
	p := &pageInfo{Links: make(map[string]int)}

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Request.AbsoluteURL(e.Attr("href"))
		if link != "" {
			p.Links[link]++
		}
		e.Request.Visit(link)
	})

	c.OnResponse(func(r *colly.Response) {
		link := r.Request.URL
		log.Infoln("response recieved", link, r.StatusCode)
		p.StatusCode = r.StatusCode
		p.End = time.Now()
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Infoln("error:", r.StatusCode, err)
		p.StatusCode = r.StatusCode
		p.End = time.Now()
	})

	p.Start = time.Now()
	c.Visit(url)

	b, err := json.Marshal(p)
	if err != nil {
		log.Errorln("failed to serialize response:", err)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.Write(b)
}
