package main

import "time"

type CrawlJob struct {
	Start  time.Time
	Finish time.Time
	*Page
	Err error
}

// NewCrawlJob will create a new job that will cause a crawl to happen
func NewCrawlJob(page *Page) (j *CrawlJob) {
	j = &CrawlJob{Page: page}
	return j
}
