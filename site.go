package main

type SiteMap map[string]*Site

type Site struct {
	Baseurl string
	PageMap
}
