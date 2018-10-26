package moni

import (
	"encoding/json"
	"testing"
)

func TestCrawlHandler(t *testing.T) {
	url := "/crawl/example.com"

	// Get the response from the handler so it can be varified
	resp := ServiceLoopback(CrawlHandler, "POST", url)
	if resp == nil {
		t.Error("CrawlHandler test failed to get a response")
	}
	body := GetBody(resp)
	if body == nil {
		t.Errorf("Crawl handler failed to read the body")
	}
	body = body
	ctype := resp.Header.Get("Content-Type")
	if ctype != "application/json" {
		t.Errorf("expected content type (application/json) got (%s)", ctype)
	}
}

func TestCrawlListHandler(t *testing.T) {
	url := "/crawlids"
	resp := ServiceLoopback(CrawlListHandler, "get", url)
	body := GetBody(resp)
	if body == nil {
		t.Errorf("Crawl list handler failed to read the body")
	}
	ctype := resp.Header.Get("Content-Type")
	if ctype != "application/json" {
		t.Errorf("expected content type (application/json) got (%s)", ctype)
	}
}

func TestCrawlIdHandler(t *testing.T) {
	url := "http://localhost:8888/store/crawl-example-com.json"
	resp := ServiceLoopback(CrawlIdHandler, "get", url)
	body := GetBody(resp)
	if body == nil {
		t.Errorf("crawl list read body failed")
	}
	ctype := resp.Header.Get("Content-Type")
	if ctype != "application/json" {
		t.Errorf("expected content type (application/json) got (%s)", ctype)
	}

	var page Page
	err := json.Unmarshal(body, &page)
	if err != nil {
		t.Errorf("failed unmarshallng body %v ", err)
	}

	if page.URL != "gardenpassages.com" {
		t.Errorf("expected (gardenpassages.com) got (%s)", page.URL)
	}
}
