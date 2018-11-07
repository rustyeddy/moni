package moni

import (
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

// Register the page routes
func registerPages(r *mux.Router) {
	r.HandleFunc("/page", PageListHandler).Methods("GET")
	r.HandleFunc("/page/", PageListHandler).Methods("GET")
	r.HandleFunc("/page/{url}", PageIdHandler).Methods("GET", "POST", "DELETE")
}

func PageListHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, app.Pagemap)
}

func PageIdHandler(w http.ResponseWriter, r *http.Request) {

	// Get a couple vars ready for later
	var err error

	// Get the url from the request and extract the storage
	// index (name) of the corresponding object.
	url := urlFromRequest(r)
	name := NameFromURL(url)

	if page := app.Pagemap.Get(url); page != nil {
		switch r.Method {
		case "GET":
			writeJSON(w, page)
			return
		case "PUT", "POST":
			log.Infoln("overwriting ", name)
		case "DELETE":
			delete(app.Pagemap, name)
		}
	} else {
		switch r.Method {
		case "GET", "DELETE":
			// Nothing to get or delete
			err = errors.New("object not found " + name)
			JSONError(w, err)
		case "PUT", "POST":
			StorePage(page)
			writeJSON(w, `{"msg": "done"}`)
		}
	}
}
