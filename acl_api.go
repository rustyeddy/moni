package moni

import "net/http"

// ACLHandler will respond to ACL requests
func ACLHandler(w http.ResponseWriter, r *http.Request) {
	acl := Crawler.AccessList
	if acl.IfNilError(acl, "crawler handler") {
		w.Write([]byte("internal error"))
		return
	}

	acl.Infoln("ACLHandler returning access list")
	writeJSON(w, acl)
}
