package main

/*
 * Moni is my website monitoring tool
 */

import (
	"flag"
	"fmt"

	"github.com/gocolly/colly"
	"github.com/rustyeddy/store"
)

// Configuration manages all variables and parameters for a given run of moni.
type Configuration struct {
	ConfigFile string
}

var (
	storage *store.FileStore
	config  Configuration

	pages []string
)

func init() {
	var err error

	flag.StringVar(&config.ConfigFile, "config", "moni.json", "Moni config file")
	storage, err = store.UseFileStore(".")
	if err != nil {
		panic(err)
	}
}

func main() {
	flag.Parse()
	c := colly.NewCollector()

	// Find and visit all links
	c.OnHTML("a", func(e *colly.HTMLElement) {
		thing := e.Attr("href")

		page = filterPage(page)
		if page == nil {
			continue
		}
		e.Request.Visit(thing)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.Visit("http://go-colly.org/")
	storage.Save("config.json", config)
}

func filterPage(page string) (page *Page) {

}
