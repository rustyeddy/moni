package moni

import "github.com/sirupsen/logrus"

// AccessList is a list of domains or url paths that are either
// allowed or denied.  AccessLists may not have both allow and deny
// entries, hence the list is one or the other.
type AccessList struct {
	Allowed     map[string]int
	Rejected    map[string]int
	Unsupported map[string]int

	*logrus.Entry `json:"-"`
}

func initACL() (acl *AccessList) {
	// ReadSaved acls if any
	if acl = ReadACL(); acl == nil {
		acl = NewACL()
	}
	return acl
}

// ACL returns the accessList.  If the accessList does not exist
// it will be created prior to return
func NewACL() *AccessList {
	acl := &AccessList{
		Allowed:     make(map[string]int),
		Rejected:    make(map[string]int),
		Unsupported: make(map[string]int),
	}
	acl.Entry = acl.logEntry()
	return acl
}

func (acl *AccessList) logEntry() *logrus.Entry {
	// straigh logrus
	flds := logrus.Fields{
		"Allowed":     len(acl.Allowed),
		"Rejected":    len(acl.Rejected),
		"Unsupported": len(acl.Unsupported),
	}
	return logrus.WithFields(flds)
}

// AllowHost will naively take only the host, ignoring port,
// and other fields to just the host.
func (acl *AccessList) AddHost(h string) {
	if host := GetHostname(h); host != "" {
		acl.Allowed[host]++
		acl.Debugln("added host ", host, " to Allowed list")
	} else {
		acl.Allowed[h] = 1
		acl.Errorln("failed to add host", host, "allowed list")
	}
}

// RemoveHost will naively take only the host, ignoring port,
// and other fields to just the host.
func (acl *AccessList) RemoveHost(h string) {
	if host := GetHostname(h); host != "" {
		delete(acl.Allowed, host)
		acl.Debugln("removed host ", host, " from Allowed list")
	} else {
		acl.Errorf("Host %s is not in ACL failed to remove", h)
	}
}

// Reject takes the host name and creates an acl entry.
// And naively ignores things like scheme and port, etc.
func (acl *AccessList) AddReject(h string) {
	if host := GetHostname(h); host != "" {
		acl.Rejected[host]++
		acl.Debugln("RejectHost ", h)
	} else {
		acl.Rejected[h] = 1
		acl.Errorln("RejectHost failed to get host for ", h)
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

// ReadACL will attempt to read the acl file (/srv/moni/acl.json)
// by default.  If it fails, it will complain then allow you to
// carry on.
func ReadACL() (acl *AccessList) {
	st := UseStore(config.Storedir)
	IfNilFatal(st)

	acl = new(AccessList)
	err := st.Get("acl.json", acl)
	IfErrorWarning(err, "reading acl.json")
	acl.Entry = acl.logEntry()
	return acl
}

func SaveACL() {
	st := UseStore(config.Storedir)
	IfNilFatal(st)

	err := st.Put("acl.json", acl)
	IfErrorFatal(err)
}
