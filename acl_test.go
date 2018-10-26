package moni

import (
	"io/ioutil"
	"net/http"
	"testing"

	log "github.com/sirupsen/logrus"
)

func TestACL(t *testing.T) {

	acl := NewACL()
	acl.AllowHost("http://allowme.com")
	acl.RejectHost("http://rejectme.com")

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

func TestGetACLHandler(t *testing.T) {
	// By creating an HTTP test client and receiver we can mock up URL
	// requests and pass the to the server which then passes the URL to
	// router (mux) and back to the callback (handler).  Beautiful.
	resp := ServiceLoopback(ACLHandler, "get", "/acl")
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
}
