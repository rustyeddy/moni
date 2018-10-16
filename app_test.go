package main

import (
	"net/http"
	"testing"
)

func TestAppHandler(t *testing.T) {
	resp := ServiceTester(t, AppHandler, "/")
	if resp == nil {
		t.Error("Failed appHandler /")
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code (20x) got (%d) ", resp.StatusCode)
	}

	ct := resp.Header.Get("Content-Type")
	if ct != "text/html" {
		t.Errorf("Expected content type (text/html) got (%s) ", ct)
	}
}
