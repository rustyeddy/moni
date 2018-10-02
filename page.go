package inv

import "time"

// ===================================================================
type Page struct {
	URL     string
	Content []byte
	Links   map[string]int

	Crawled    bool
	StatusCode int
	Start      time.Time
	End        time.Time
}

// ********************************************************************
type PageMap map[string]*Page

func (pm PageMap) Get(url string) (p *Page) {
	if p, e := Visited[url]; e {
		return p
	}
	return nil
}

func (pm PageMap) Set(url string, p *Page) {
	Visited[url] = p
}
