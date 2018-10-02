package inv

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

var (
	CrawlDepth int = 1
)

func HandleCrawl(w http.ResponseWriter, r *http.Request) {

	// Using Gorilla mux /walk/{url}
	vars := mux.Vars(r)
	url := vars["url"]

	// TODO: check for "http://" prefix
	if !strings.HasPrefix(url, "http") {
		url = "http://" + url
	}

	/*
		go Filter(filterCh <-chan)
		go Crawl(filterCh, storeCh)
		go Store(storeCh)
		filterCh <- url
	*/

	// Make this a go routine
	Crawl(url)

	// Done crawling make json object to send back
	jbytes, err := json.Marshal(Visited)
	if err != nil {
		log.Errorln("failed to serialize response:", err)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(jbytes)

}

func Crawl(url string) (p *Page) {
	log.Infoln("crawling", url)

	// Create the collector and go get shit!
	c := colly.NewCollector(
		colly.MaxDepth(CrawlDepth),
	)

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Request.AbsoluteURL(e.Attr("href"))
		if link == "" {
			// do we even want to record these?
			log.Infoln("  link is nil %+v\n", e)
			//return
		}
		p.Links[link]++
		e.Request.Visit(link)
	})

	c.OnResponse(func(r *colly.Response) {
		link := r.Request.URL
		log.Infoln("response recieved", link, r.StatusCode)
		p.StatusCode = r.StatusCode
		p.End = time.Now()
		p.Crawled = true
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Infoln("error:", r.StatusCode, err)
		p.StatusCode = r.StatusCode
		p.End = time.Now()
	})

	// Get the page we'll use for this walk
	if p = Visited.Get(url); p == nil {
		p = &Page{
			URL:   url,
			Links: make(map[string]int),
		}
		Visited[url] = p
	}
	p.Start = time.Now()
	c.Visit(url)
	return p
}
