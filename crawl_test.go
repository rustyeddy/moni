package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCrawlHandler(t *testing.T) {

	// Craft up a request with the URL we want to test
	req := httptest.NewRequest("GET", "http://example.com:8888/crawl/example.com", nil)
	w := httptest.NewRecorder()

	// Now crawl
	CrawlHandler(w, req)

	// Look at the response
	resp := w.Result()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("failed to read body %v", err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("expected (200) got (%d)", resp.StatusCode)
	}

	ctype := resp.Header.Get("Content-Type")
	if ctype != "application/json" {
		t.Errorf("expected content type (application/json) got (%s)", ctype)
	}

	data := string(body)
	if data != "foo" {
		t.Errorf("expected body (%s) got (%s)", "foo", data)
	}
}

// handler keeps failing ???
func TestACLHandler(t *testing.T) {
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("GET", "/acl", nil)
	if err != nil {
		t.Errorf("failed to create new test %v", err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ACLHandler)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	bod := rr.Body.String()
	fmt.Printf("BODY STRING %s", bod)
	expected := `{"alive": true}`
	if bod != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			bod, expected)
	}
}
