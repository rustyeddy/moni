package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// ServiceTester will pass the given handler and url to a special
// http test client and and "server", bypassing the network.  The
// Request, URL and variables are sanitized and prepared before
// calling the Handler function (just like a real server).  The request
// returned by the Handler is returned to the calling test for
// examination.  The calling test has the context required to determine
// the passability of the test.
func ServiceTester(t *testing.T, h http.HandlerFunc, url string) *http.Response {

	// Craft up a request with the URL we want to test
	req := httptest.NewRequest("GET", url, nil)
	w := httptest.NewRecorder()

	// Do not give the writer and request to the handler directly.  The args
	// will not have been processed.  Register it as a handler the let mux
	// parse our the args and setup other important things, then let it
	// call the handler itself.
	//handler := http.HandlerFunc(h)
	handler := httpServer().Handler

	// This will cause the actual crawling, CrawlHandler will be called
	// with all the appropriate header and argument processing.
	handler.ServeHTTP(w, req)

	// Look at the response
	resp := w.Result()
	return resp
}

func TestCrawlHandler(t *testing.T) {
	url := "/crawl/example.com"

	// Get the response from the handler so it can be varified
	resp := ServiceTester(t, CrawlHandler, url)
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
	url := "/crawls"
	resp := ServiceTester(t, CrawlListHandler, url)
	body := GetBody(resp)
	if body == nil {
		t.Errorf("Crawl list handler failed to read the body")
	}
	ctype := resp.Header.Get("Content-Type")
	if ctype != "application/json" {
		t.Errorf("expected content type (application/json) got (%s)", ctype)
	}
}
