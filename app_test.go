package moni

import (
	"net/http"
	"testing"
)

var (
	tapp *App
)

func init() {
	tapp = NewApp(&DefaultConfig)
}

// TestAppHandler checks what we get with /
func TestAppRootHandler(t *testing.T) {
	resp := ServiceLoopback(AppHandler, "get", "/")
	if resp == nil {
		t.Error("Failed appHandler /")
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code (20x) got (%d) ", resp.StatusCode)
	}

	ct := resp.Header.Get("Content-Type")
	if ct == "" || len(ct) < 9 {
		t.Errorf("Expected content type (text/html) got (%s) ", ct)
	}

	// TODO ~ Check the body make sure the content is corrent
}
