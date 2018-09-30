package main

import (
	"encoding/json"
	"net/http"

	"github.com/gocolly/colly"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type pageInfo struct {
	StatusCode int
	Links      map[string]int
}

func handleWalk(w http.ResponseWriter, r *http.Request) {

	// Using Gorilla mux /walk/{url}
	vars := mux.Vars(r)
	url := vars["url"]

	// This is standard {url}
	// url := r.URL.Query().Get("url")
	// if url == "" {
	// 	log.Errorln("missing URL")
	// 	return
	// }

	url = "http://" + url
	log.Infoln("visiting", url)

	c := colly.NewCollector(
		colly.MaxDepth(3),
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
		log.Infoln("response recieved", r.Request.URL, r.StatusCode)
		p.StatusCode = r.StatusCode
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Infoln("error:", r.StatusCode, err)
		p.StatusCode = r.StatusCode
	})

	c.Visit(url)

	b, err := json.Marshal(p)
	if err != nil {
		log.Errorln("failed to serialize response:", err)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.Write(b)
}
