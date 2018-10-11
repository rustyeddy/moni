package main

import (
	"fmt"
	"time"
)

// ===================================================================
type Page struct {
	URL string

	Content []byte
	Links   map[string]*Page
	Ignored map[string]int

	Crawled    bool
	StatusCode int

	Start time.Time
	End   time.Time
}

// ********************************************************************
type Pagemap map[string]*Page

func GetPage(url string) (pi *Page) {
	var ex bool
	if pi, ex = Pages[url]; !ex {
		pi = &Page{
			URL:     url,
			Links:   make(map[string]*Page),
			Ignored: make(map[string]int),
		}
		Pages[url] = pi
	}
	fmt.Printf(" returning page %+v\n", pi)
	return pi
}

func (pm Pagemap) Get(url string) (p *Page) {
	if p, e := Visited[url]; e {
		return p
	}
	return nil
}

func (pm Pagemap) Exists(url string) bool {
	if p := pm.Get(url); p != nil {
		return true
	}
	return false
}

func (pm Pagemap) Set(url string, p *Page) {
	Visited[url] = p
}
