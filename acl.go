package main

// AccessList is a list of domains or url paths that are either
// allowed or denied.  AccessLists may not have both allow and deny
// entries, hence the list is one or the other.
type AccessList map[string]bool

func (acl AccessList) Add(url string, allow bool) {
	acl[url] = allow
}

func (acl AccessList) Allow(url string) bool {
	if ac, e := acl[url]; e {
		return ac
	}
	return false
}

func (acl *AccessList) Reject(url string) bool {
	return !acl.Allow(url)
}
