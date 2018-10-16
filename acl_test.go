package main

import (
	"io/ioutil"
	"net/http"
	"testing"

	log "github.com/sirupsen/logrus"
)

func TestNormalizeURL(t *testing.T) {
	var tsts = []struct {
		host     string
		expected string
	}{
		{"amazon.com", "http://amazon.com"},
		{"//clowdops.net:4040", "http://clowdops.net:4040"},
		{"http://rustyeddy.com", "http://rustyeddy.com"},
		{"//example.com:300", "http://example.com:300"},
		{"//john.bozo.nono.com:3000", "http://john.bozo.nono.com:3000"},
		{"tel:phonenumber", ""},
		{"", ""},
	}

	for _, tst := range tsts {
		ustr, err := NormalizeURL(tst.host)
		if err != nil && tst.expected != "" {
			// Here if we have an error and did not expect one
			t.Errorf("Normalize URL failed %v", err)
		}
		if ustr != tst.expected {
			t.Errorf("Normalize URL (%s) expected (%s) got (%s)", tst.host, tst.expected, ustr)
		}
	}
}

func TestACL(t *testing.T) {
	acl := AccessList{
		Allowed:     make(map[string]int),
		Rejected:    make(map[string]int),
		Unsupported: make(map[string]int),
	}

	acl.AllowHost("allowme.com")
	acl.RejectHost("rejectme.com")

	var tsts = []struct {
		host     string
		expected bool
	}{
		{"allowme.com", true},   // allowed
		{"rejectmd.com", false}, // explicit reject
		{"unknown.com", false},  // reject by default
	}

	for _, tst := range tsts {
		allowed := acl.IsAllowed(tst.host)
		if allowed != tst.expected {
			log.Errorf("acl allow(%s) expected (%t) got (%t) ", tst.host, tst.expected, allowed)
		}
	}
}

func TestACLHandler(t *testing.T) {
	// By creating an HTTP test client and receiver we can mock up URL
	// requests and pass the to the server which then passes the URL to
	// router (mux) and back to the callback (handler).  Beautiful.
	resp := ServiceTester(t, ACLHandler, "/acl")
	if resp == nil {
		t.Error("CrawlHandler test failed to get a response")
		return
	}

	// Check the status code is what we expect.
	if resp.StatusCode != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			resp.StatusCode, http.StatusOK)
	}

	// Check the response body is what we expect.
	_, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("ready respone body %v", err)
	}
	/*
		expected := `{"alive": true}`
		if string(bod) != expected {
			t.Errorf("handler returned unexpected body: got %v want %v",
				bod, expected)
		}
	*/
}
