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
	writeJSON(w, app.AccessList)
}

func ACLHandler(w http.ResponseWriter, r *http.Request) {
	var url string
	if url = urlFromRequest(r); url == "" {
		JSONError(w, errors.New("url From Request failed "+url))
		return
	}

	switch r.Method {
	case "GET":
		if app.IsAllowed(url) {
			writeJSON(w, true)
		} else {
			writeJSON(w, struct{}{})
		}

	case "PUT", "POST":
		app.AddHost(url)
		writeJSON(w, struct{ url string }{url})

	case "DELETE":
		if app.IsAllowed(url) {
			app.RemoveHost(url)
		}

	default:
		log.Errorf("method %s unknown", r.Method)
	}
	SaveACL()
}
