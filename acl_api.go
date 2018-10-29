package moni

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func init() {

}

func registerACLHandler(r *mux.Router) {
	r.HandleFunc("/acl", ACLHandler) // Display ACLs
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
