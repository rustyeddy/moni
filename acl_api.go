package moni

import (
	"fmt"
	"net/http"
)

func init() {

}

// ACLHandler will respond to ACL requests
func ACLHandler(w http.ResponseWriter, r *http.Request) {
	acl := Crawler.AccessList

	if IfNilError(acl, "crawler handler") {
		JSONError(w, fmt.Errorf("Expected (acl) got () "))
		return
	}

	acl.Infoln("ACLHandler returning access list")
	writeJSON(w, acl)
}
