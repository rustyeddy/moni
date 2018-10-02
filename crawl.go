package inventory

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gocolly/colly"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

var (
	CrawlDepth int = 1
)

func HandleCrawl(w http.ResponseWriter, r *http.Request) {

	fmt.Println("  Handle Crawl ")

	// Using Gorilla mux /walk/{url}
	vars := mux.Vars(r)
	url := vars["url"]
	fmt.Println("  Handle Crawl ", url)

	// TODO: check for "http://" prefix
	url = "http://" + url
	p := Crawl(url)
	b, err := json.Marshal(p)
	if err != nil {
		log.Errorln("failed to serialize response:", err)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(b)

}

func Crawl(url string) (p *PageInfo) {
	log.Infoln("visiting", url)

	// Create the collector and go get shit!
	c := colly.NewCollector(
		colly.MaxDepth(CrawlDepth),
	)
	p = &PageInfo{Links: make(map[string]int)}

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
	return p
}
