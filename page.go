package inv

import "time"

type PageMap map[string]*Page

var (
	Visited PageMap = make(PageMap)
)

func (pm PageMap) Get(url string) (p *Page) {
	if p, e := Visited[url]; e {
		return p
	}
	return nil
}

func (pm PageMap) Set(url string, p *Page) {
	Visited[url] = p
}

type Page struct {
	URL        string
	StatusCode int
	Links      map[string]int
	Crawled    bool
	Start      time.Time
	End        time.Time
}
