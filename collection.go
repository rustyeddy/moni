package moni

import (
	"fmt"
	"time"

	"github.com/gocolly/colly"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
	log "github.com/sirupsen/logrus"
)

type RoundTrip struct {
	URL    string
	Start  time.Time
	Finish time.Time
	Err    error
}

func NewRoundTrip(url string) (rtt *RoundTrip) {
	rtt = &RoundTrip{URL: url}
	return rtt
}

// Collection (not to be confused with a Mongo Collection) manages
// groups of related crawls.
type Collection struct {
	Name string

	Pagemap // A collection of page
	Links   map[string]*RoundTrip

	*mongo.Collection
	*bson.Document

	*colly.Collector
	*log.Logger

	Callbacks []func() // This function will register all callbacks
	// creating a function allows us to change the behavior by constructing
	// different sets of callbacks to suit the application

	FilteredLinks int
	urlch         chan string
}

func NewCollection(cb func()) (c *Collection) {
	c = &Collection{
		Collector: colly.NewCollector(
			colly.MaxDepth(1),
			//colly.AllowedDomains("http://sierrahydrographics.com"),
			colly.DisallowedDomains("namecheap.com", "www.namecheap.com", "wordpress.org", "www.wordpress.org", "developer.wordpress.org"),
			colly.Async(true),
		),
		Logger: GetDebugLogger(),
	}

	// Limit parallelism to 2
	c.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 2})

	// Just connect all callbacks!
	c.Callbacks = append(c.Callbacks, c.LinkCrawlerCallbacks)
	c.Callbacks = append(c.Callbacks, c.OutlineCallbacks)
	return c
}

func (c *Collection) ScheduleCollection() {
	log.Debugln("Scheduling collection ", c.Name)

	// We do not have a scheduler ... yet just do it
	//StartCollection()
	for _, cb := range c.Callbacks {
		cb()
	}
}

func (c *Collection) LinkCrawlerCallbacks() {
	c.Links = make(map[string]*RoundTrip)
	c.OnRequest(func(r *colly.Request) {
		url := r.URL.String()
		c.Links[url] = NewRoundTrip(url)
	})

	c.OnError(func(r *colly.Response, err error) {
		url := r.Request.URL.String()
		c.Links[url].Err = err
	})

	c.OnResponse(func(r *colly.Response) {
		url := r.Request.URL.String()
		if url == "" {
			c.Errorln("OnResponse failed to get url")
			return
		}
		_, ex := c.Links[url]
		if !ex {
			c.Links[url] = NewRoundTrip(url)
		}
		c.Links[url].Finish = time.Now()
	})

	// ==========================================================

	// On every a element which has href attribute call callback
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		fmt.Println(" in element land ")
		if url := e.Attr("href"); url != "" {
			fmt.Println("  crawl ", url)
			c.CrawlURL(url)
		}
	})

	c.OnHTML("tr td:nth-of-type(1)", func(e *colly.HTMLElement) {
		c.Infoln("  table col1: ", e.Text)
	})

	c.OnXML("//h1", func(e *colly.XMLElement) {
		//c.Infoln("  h1: ", e.Text)
	})

	c.OnXML("//h2", func(e *colly.XMLElement) {
		//c.Infoln("  h2: ", e.Text)
	})

	c.OnXML("//h3", func(e *colly.XMLElement) {
		//c.Infoln("  h2: ", e.Text)
	})

	c.OnScraped(func(r *colly.Response) {
		//c.Infoln("Finished", r.Request.URL)
	})
}

func (c *Collection) OutlineCallbacks() {
	panic("Must implement before use")
}

// Process links will listen to the linkChannel for incoming or new links
// This process will also be responsible for filtering links
/*
func (c *Collection) ProcessLinks(ch chan string) {
	for {
		url := <-c.urlch
		c.Debugln("  incoming url! ", url)
		if link := c.FilterURL(url); link != nil {
			c.qLink(link)
		} else {
			c.FilteredLinks++
		}
	}
}
*/
func (c *Collection) CrawlURL(url string) {
	fmt.Println("Sending url to channel ... ", url)
	c.urlch <- url
	fmt.Println("Sent ...")
}
