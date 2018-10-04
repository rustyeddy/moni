package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/gorilla/mux"

	//"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

/*
  1. filterChannel <- newURL

  2. if reject(newURL) goto 1.

  3. visitChannel <- newURL

  4. Page = NewPage(newURL)

  5. RecordPage(p)  // used by filter

  6. VisitPage(p)
     - Page.Start = time.Now()
     - callback: OnRequest()

  7. RecvReply(resp)
     - callback: OnResponse
     - Page.End = time.Now() // record end of RTT

  8. Parse(resp)

  8.1 callback: OnDocument
  8.2 callback: OnElement["href"]
      - anchor, aLink = LinkFromResp(resp)
      - FilterChan <- aLink
  8.3 callback: OnHTML

  9. StorageChan <-DOM
*/

var (
	Visited    PageMap = make(PageMap)
	CrawlDepth int     = 1
)

func HandleCrawl(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	url := vars["url"]

	if !strings.HasPrefix("http", url) {
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
		//colly.MaxDepth(config.Depth),
		colly.MaxDepth(3),
	)

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Request.AbsoluteURL(e.Attr("href"))
		if link == "" {
			return
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
