package main

// Pages is the structure that maintains all pages for this
// root site
type Pages struct {
	BaseURL string
	Pages   map[string]*Page
}

type Page struct {
	URL string
}
