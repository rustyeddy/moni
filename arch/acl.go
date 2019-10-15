package moni

import (
	"encoding/json"
	"io/ioutil"

	"github.com/rustyeddy/store"
	log "github.com/sirupsen/logrus"
)

// AccessList is a list of domains or url paths that are either
// allowed or denied.  AccessLists may not have both allow and deny
// entries, hence the list is one or the other.
type AccessList struct {
	Allowed     map[string]int
	Rejected    map[string]int
	Unsupported map[string]int

	*log.Entry `json:"-"` // JSON to ignore this
}

// ACL returns a brand new accessList with all default values.
func NewACL() AccessList {
	acl := AccessList{
		Allowed:     make(map[string]int),
		Rejected:    make(map[string]int),
		Unsupported: make(map[string]int),
	}
	acl.Entry = acl.logEntry()
	return acl
}

func initACL() (acl *AccessList) {
	if acl = ReadACL(); acl == nil {
		*acl = NewACL()
	}
	return acl
}

func (acl *AccessList) logEntry() *log.Entry {
	// straigh logrus
	flds := log.Fields{
		"Allowed":     len(acl.Allowed),
		"Rejected":    len(acl.Rejected),
		"Unsupported": len(acl.Unsupported),
	}
	return log.WithFields(flds)
}

// AllowHost will naively take only the host, ignoring port,
// and other fields to just the host.
func (acl *AccessList) AddHost(h string) {
	if host := GetHostname(h); host != "" {
		if _, ex := acl.Allowed[host]; ex {
			acl.Allowed[host]++
		} else {
			acl.Allowed[host] = 1
		}
	} else {
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

// ReadACL will attempt to read the acl file (/etc/moni/acl.json)
// by default.  If it fails, it will complain then allow you to
// carry on.
func ReadACL() *AccessList {
	var (
		buf []byte
		a   AccessList
		err error
	)

	storedir := "./etc"
	st, err := store.UseFileStore(storedir)
	IfNilFatal(st)

	// Err and Panic
	IfErrorFatal(err, "Failed to open store at ")

	path := "./etc/acl.json"
	if buf, err = ioutil.ReadFile(path); err != nil {
		log.Errorf("read index %s failed %v", path, err)
		return nil
	}

	// assume a content type of JSON (use the .ext)
	if err = json.Unmarshal(buf, &a); err != nil {
		IfErrorFatal(err, "get failed marshaling json "+path)
		return nil
	}

	/*
			nacl := AccessList{
				Allowed:     make(map[string]int),
				Rejected:    make(map[string]int),
				Unsupported: make(map[string]int),
			}
		fmt.Printf("going to get acl.json")
		if err := st.Get("acl.json", &nacl); err != nil {
			log.Warnf("acl.json could not be read %v", err)
			return nil
		}
	*/
	return &a
}

// SaveACL stores the acl in our store
func SaveACL() {
	st, _ := store.UseFileStore("./etc")
	IfNilFatal(st)

	err := st.Put("acl.json", acl)
	IfErrorFatal(err)
}
