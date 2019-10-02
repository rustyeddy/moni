package moni

import (
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func registerACLHandler(r *mux.Router) {
	r.HandleFunc("/acl", ACLListHandler)
	r.HandleFunc("/acl/", ACLListHandler).Methods("GET", "DELETE")
	r.HandleFunc("/acl/{url}", ACLHandler).Methods("GET", "PUT", "POST", "DELETE")
}

// ACLHandler will respond to ACL requests
func ACLListHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, acl)
}

func ACLHandler(w http.ResponseWriter, r *http.Request) {
	var url string
	if url = urlFromRequest(r); url == "" {
		JSONError(w, errors.New("url From Request failed "+url))
		return
	}

	switch r.Method {
	case "GET":
		if acl.IsAllowed(url) {
			writeJSON(w, true)
		} else {
			writeJSON(w, struct{}{})
		}

	case "PUT", "POST":
		acl.AddHost(url)
		writeJSON(w, struct{ url string }{url})

	case "DELETE":
		if acl.IsAllowed(url) {
			acl.RemoveHost(url)
		}

	default:
		log.Errorf("method %s unknown", r.Method)
	}
}
