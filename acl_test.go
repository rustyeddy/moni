package main

import (
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
		var ustr string

		u, err := NormalizeURL(tst.host)
		if err != nil {
			ustr = ""
		} else {
			ustr = u.String()
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
