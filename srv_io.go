package moni

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// Write the response back to the caller as JSON
func writeJSON(w http.ResponseWriter, obj interface{}) {
	var jbytes []byte
	var err error

	// Set response content type to json
	if jbytes, err = json.Marshal(obj); err != nil {
		http.Error(w, err.Error(), 500) // TODO: replace 500 with http.XXX
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(jbytes)
	if err != nil {
		JSONError(w, err)
	}
}

func JSONError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	http.Error(w, err.Error(), 500)
	log.Error(err)
}

func HTMLError(w http.ResponseWriter, err string) {
	w.Header().Set("Content-Type", "text/html")
	http.Error(w, err, 500)
	log.Error(err)
}
