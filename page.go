package main

import (
	"fmt"
	"net/url"
	"time"
)

// Page represents a single web page
type Page struct {
	*Site      `json:"-"`
	url.URL    `json:"url"`
	Links      map[string][]string `json:"links"`
	StatusCode int                 `json:"statusCode"`

	*time.Ticker `json:"-"`
	TimeStamp    `json:"timestamp"`
	TimeStamps   []TimeStamp `json:"timestamps"`
}

func (p *Page) PageString() (s string) {
	s = fmt.Sprintf("%-40s: links: %-4d resp: %-10v ", p.URL.String(), len(p.Links), p.Elapsed)
	if tslen := len(p.TimeStamps); tslen > 0 {
		if tslen > 4 {
			tslen -= 4
		}
		s += fmt.Sprintf("\tlast elasped: %v", p.TimeStamps[tslen:])
	}
	s += "\n"

	if config.Verbose {
		for l, _ := range p.Links {
			s += fmt.Sprintf("\t%s\n", l)
		}
	}
	return s
}
