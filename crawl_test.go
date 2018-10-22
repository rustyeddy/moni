package moni

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	log "github.com/sirupsen/logrus"
)

// ServiceTester will pass the given handler and url to a special
// http test client and and "server", bypassing the network.  The
// Request, URL and variables are sanitized and prepared before
// calling the Handler function (just like a real server).  The request
// returned by the Handler is returned to the calling test for
// examination.  The calling test has the context required to determine
// the passability of the test.
func ServiceTester(t *testing.T, h http.HandlerFunc, verb string, url string) *http.Response {

	// Craft up a request with the URL we want to test
	req := httptest.NewRequest(verb, url, nil)
	w := httptest.NewRecorder()

	// Do not give the writer and request to the handler directly.  The args
	// will not have been processed.  Register it as a handler the let mux
	// parse our the args and setup other important things, then let it
	// call the handler itself.
	//handler := http.HandlerFunc(h)
	r := Server().Handler

	// This will cause the actual crawling, CrawlHandler will be called
	// with all the appropriate header and argument processing.
	r.ServeHTTP(w, req)

	// Look at the response
	resp := w.Result()
	if resp == nil {
		log.Errorln("failed to get a response")
	}
	return resp
}

func TestCrawlHandler(t *testing.T) {
	url := "/crawl/example.com"

	// Get the response from the handler so it can be varified
	resp := ServiceTester(t, CrawlHandler, "POST", url)
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
	resp := ServiceTester(t, CrawlListHandler, "get", url)
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
	url := "var/tstore/example-com.json"
	resp := ServiceTester(t, CrawlIdHandler, "get", url)
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
