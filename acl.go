package main

import (
	"net/http"
	"net/url"

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

func GetHostname(h string) (host string) {
	var err error
	var u *url.URL

	if u, err = url.Parse(h); err != nil {
		log.Errorf("failed to parse hostname", err)
		return ""
	}
	return u.Hostname()
}

// AllowHost will naively take only the host, ignoring port,
// and other fields to just the host.
func (acl *AccessList) AllowHost(h string) {
	if host := GetHostname(h); host != "" {
		acl.Allowed[host]++
		log.Debugln("added host ", host, " to Allowed list")
	} else {
		log.Errorln("failed ot add host", host, "allowed list")
	}
}

// Reject takes the host name and creates an acl entry.
// And naively ignores things like scheme and port, etc.
func (acl *AccessList) RejectHost(h string) {
	if host := GetHostname(h); host != "" {
		acl.Rejected[host]++
		log.Debugln("RejectHost ", h)
	} else {
		log.Errorln("RejectHost failed to get host for ", h)
	}
	return
}

// IsAllowed matches the url against the acl to determine
// if this site (or url) is allowed to be crawled.  IsAllowed
// can assume the  url has been normalized
func (acl *AccessList) IsAllowed(urlstr string) (allow bool) {
	host := GetHostname(urlstr)
	if _, ex := acl.Allowed[host]; ex {
		acl.Allowed[host]++
		return true
	}
	acl.Rejected[host]++
	return false
}

// ACLHandler will respond to ACL requests
func ACLHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, &ACL)
}
