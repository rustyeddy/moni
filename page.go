package main

import "net/url"

type Page struct {
	*url.URL
	StatusCode int
	Links      map[string]int
}

// NewPage returns a new page structure
func NewPage(urlstr string) (p *Page) {
	u, err := url.Parse(urlstr)
	ErrPanic(err)

	p = &Page{
		URL:        u,
		StatusCode: 0,
		Links:      make(map[string]int),
	}
	return p
}

func ErrPanic(err error) {
	if err != nil {
		panic(err)
	}
}
