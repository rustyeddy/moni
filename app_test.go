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

// TestStaticFiles will verify we can access static files such as css and
// javascript along with our index files, etc..
func TestStaticFiles(t *testing.T) {
	var resp *http.Response
	if resp = ServiceLoopback(AppHandler, "get", "/css/app.css"); resp == nil {
		t.Error("expected /static/css/app.css got (nil) ")
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code (20x) got (%d) ", resp.StatusCode)
	}
	ct := resp.Header.Get("Content-Type")
	if ct == "" {
		t.Errorf("Expected content type (text/html) got (%s) ", ct)
	}
	if len(ct) < len("text/css") {
		t.Errorf("Expected content type (text/html) got (%s) ", ct)
	}
}
