package moni

import (
	"encoding/json"
	"testing"
)

func TestSiteListHandler(t *testing.T) {

	url := "/site/"
	resp := ServiceLoopback(SiteListHandler, "get", url)

	body := GetBody(resp)
	if body == nil {
		t.Errorf("crawl list read body failed")
	}
	ctype := resp.Header.Get("Content-Type")
	if ctype != "application/json" {
		t.Errorf("expected content type (application/json) got (%s)", ctype)
	}

	var sites []*Site
	err := json.Unmarshal(body, &sites)
	if err != nil {
		t.Errorf("failed unmarshallng sites %v ", err)
	}

	if len(sites) < 0 {
		t.Errorf("failed should have more sites")
	}
}

func TestSiteIdHandler(t *testing.T) {
	url := "/site/rustyeddy.com"
	var body []byte

	resp := ServiceLoopback(SiteIdHandler, "get", url)
	if body = GetBody(resp); body == nil {
		t.Errorf("failed to get body from response")
	}
	//fmt.Printf("site %s %+v", resp.Request.URL, body)
}
