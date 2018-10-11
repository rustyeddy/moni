package main

import (
	"encoding/json"
	"fmt"
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

func (acl *AccessList) AllowHost(ustr string) {
	u, _ := NormalizeURL(ustr)
	if u == nil {
		log.Errorln("AllowHost normalizedURL failed to")
		return
	}
	hname := u.Hostname()
	acl.Allowed[hname]++
}

func (acl *AccessList) RejectHost(ustr string) {
	u, _ := NormalizeURL(ustr)
	if ustr == "" {
		log.Errorln("RejectHost NormalizedURL failed")
		return
	}
	hname := u.Hostname()
	acl.Rejected[hname]++
}

func (acl *AccessList) IsAllowed(ustr string) (allow bool) {
	u, err := url.Parse(ustr)
	if err != nil {
		return false
	}
	hname := u.Hostname()
	if hname == "" {
		acl.Rejected[hname]++
		return false
	}

	if _, ex := acl.Allowed[hname]; ex {
		acl.Allowed[hname]++
		return true
	}

	if _, ex := acl.Rejected[hname]; ex {
		acl.Rejected[hname]++
	} else {
		acl.Rejected[hname]++
	}
	return false
}

func ACLHandler(w http.ResponseWriter, r *http.Request) {

	jbytes, err := json.Marshal(&ACL)
	if err != nil {
		log.Errorf("failed to marshal ACL %v", err)
		fmt.Fprintln(w, "internal error")
	} else {
		w.Write(jbytes)
	}

	// Now lets try to store this thing
	_, err = Storage.StoreObject("acl", &ACL)
	if err != nil {
		log.Errorf("ACL Handler ~ failed to store ACL")
	}
}
