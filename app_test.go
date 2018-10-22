package main

import (
	"net/http"
	"testing"
)

// TestAppHandler checks what we get with /
func TestAppRootHandler(t *testing.T) {
	resp := ServiceTester(t, AppHandler, "get", "/")
	if resp == nil {
		t.Error("Failed appHandler /")
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code (20x) got (%d) ", resp.StatusCode)
	}

	ct := resp.Header.Get("Content-Type")
	if ct[0:9] != "text/html" {
		t.Errorf("Expected content type (text/html) got (%s) ", ct)
	}
}

// TestStaticFiles will verify we can access static files such as css and
// javascript along with our index files, etc..
func TestStaticFiles(t *testing.T) {
	var resp *http.Response
	if resp = ServiceTester(t, AppHandler, "get", "/css/app.css"); resp == nil {
		t.Error("expected /css/app.css got (nil) ")
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code (20x) got (%d) ", resp.StatusCode)
	}

	ct := resp.Header.Get("Content-Type")
	if ct[0:9] != "text/html" {
		t.Errorf("Expected content type (text/html) got (%s) ", ct)
	}
}
