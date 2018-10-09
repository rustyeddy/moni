package main

import "time"

// ===================================================================
type PageInfo struct {
	URL     string
	Content []byte
	Links   map[string]int
	Ignored map[string]int

	Crawled    bool
	StatusCode int

	Start time.Time
	End   time.Time
}

// ********************************************************************
type PageInfomap map[string]*PageInfo

func GetPageInfo(url string) (pi *PageInfo) {
	if pi, ex := Pages[url]; !ex {
		pi = &PageInfo{
			URL:     url,
			Links:   make(map[string]int),
			Ignored: make(map[string]int),
		}
		Pages[url] = pi
	}
	return pi
}

func (pm PageInfomap) Get(url string) (p *PageInfo) {
	if p, e := Visited[url]; e {
		return p
	}
	return nil
}

func (pm PageInfomap) Exists(url string) bool {
	if p := pm.Get(url); p != nil {
		return true
	}
	return false
}

func (pm PageInfomap) Set(url string, p *PageInfo) {
	Visited[url] = p
}
