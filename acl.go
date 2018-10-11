package main

import log "github.com/sirupsen/logrus"

// AccessList is a list of domains or url paths that are either
// allowed or denied.  AccessLists may not have both allow and deny
// entries, hence the list is one or the other.
type AccessList struct {
	Allowed     map[string]int
	Rejected    map[string]int
	Unsupported map[string]int
}

func (acl *AccessList) AllowHost(ustr string) {
	if _, ustr = NormalizeURL(ustr); ustr == "" {
		log.Errorln("AllowHost normalizedURL failed to")
		return
	}
	acl.Allowed[ustr]++
}

func (acl *AccessList) RejectHost(ustr string) {
	if _, ustr = NormalizeURL(ustr); ustr == "" {
		log.Errorln("RejectHost NormalizedURL failed")
		return
	}
	return
}

func (acl *AccessList) IsAllowed(ustr string) (allow bool) {
	_, nustr := NormalizeURL(ustr)
	if nustr == "" {
		log.Infoln("AccessList failed to normalize ", ustr)
		return false
	}
	if _, ex := acl.Allowed[nustr]; ex {
		acl.Allowed[nustr]++
		return true
	}

	if _, ex := acl.Rejected[nustr]; ex {
		acl.Rejected[nustr]++
	} else {
		acl.Rejected[nustr]++
	}
	return false
}
