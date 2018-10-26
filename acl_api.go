package moni

import "net/http"

// ACLHandler will respond to ACL requests
func ACLHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, accessList)
}
