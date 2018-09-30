package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gocolly/colly"
	log "github.com/sirupsen/logrus"
)

func handleWalk(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query().Get("url")
	if url == "" {
		log.Errorln("missing URL")
		return
	}
	url = "http://" + url
	fmt.Println("visiting", url)

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
		log.Infoln("response recieved", r.StatusCode)
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
