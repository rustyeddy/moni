package main

import (
	"testing"

	log "github.com/sirupsen/logrus"
)

func TestACL(t *testing.T) {
	var tsts = []struct {
		host     string
		expected bool
	}{
		{"amazon.com", false},
		{"clowdops.net", true},
		{"http://rustyeddy.com", true},
		{"http://example.com", true},
		{"tel:phonenumber", false},
	}

	acl := AccessList{
		Allowed:     make(map[string]int),
		Rejected:    make(map[string]int),
		Unsupported: make(map[string]int),
	}

	acl.AllowHost("oclowvision.com")
	acl.AllowHost("clowdops.net")
	acl.AllowHost("example.com")
	acl.AllowHost("rustyeddy.com")

	acl.RejectHost("amazon.com")
	acl.RejectHost("walmart.com")
	acl.RejectHost("ebay.com")

	for _, tst := range tsts {
		allowed := acl.IsAllowed(tst.host)
		if allowed != tst.expected {
			log.Errorf("acl allow(%s) expected (%t) got (%t) ", tst.host, tst.expected, allowed)
		}
	}

}
