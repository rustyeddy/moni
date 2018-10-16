package main

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

// AccessList is a list of domains or url paths that are either
// allowed or denied.  AccessLists may not have both allow and deny
// entries, hence the list is one or the other.
type AccessList struct {
	Allowed     map[string]int
	Rejected    map[string]int
	Unsupported map[string]int
}

// AllowHost will naively take only the host, ignoring port,
// and other fields to just the host.
func (acl *AccessList) AllowHost(h string) {
	if h == "" {
		log.Errorln("AllowHost normalizedURL failed to")
		return
	}
	acl.Allowed[h]++
}

// Reject takes the host name and creates an acl entry.
// And naively ignores things like scheme and port, etc.
func (acl *AccessList) RejectHost(h string) {
	if h == "" {
		log.Errorln("RejectHost NormalizedURL failed")
		return
	}
	acl.Rejected[h]++
}

// IsAllowed matches the url against the acl to determine
// if this site (or url) is allowed to be crawled.  IsAllowed
// can assume the  url has been normalized
func (acl *AccessList) IsAllowed(hp string) (allow bool) {
	if _, ex := acl.Allowed[hp]; ex {
		acl.Allowed[hp]++
		return true
	}

	if _, ex := acl.Rejected[hp]; ex {
		acl.Allowed[hp]++
	} else {
		acl.Rejected[hp]++
	}
	return false
}

// ACLHandler will respond to ACL requests
func ACLHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, &ACL)
}
