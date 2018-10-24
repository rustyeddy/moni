package moni

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Get storage will return import information about the storage
// we are using.
func registerStorage(r *mux.Router) {
	r.HandleFunc("/storage", StorageHandler).Methods("GET")
}

// StorageHandler manages requests from a client
func StorageHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, GetStorage())
}
