package main

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// Write the response as HTML
/*
func writeHTML(w http.ResponseWriter, title string, body string) {
	if Config.Hasfiles {
		t, err := template.ParseFiles("tmpl/index.html")
		if err != nil {
			HTMLError(w, err.Error())
			return
		}
		p := NewPage(title, body)
		if err = t.Execute(w, p); err != nil {
			HTMLError(w, err.Error())
		}
	}
}
*/

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
	http.Error(w, err, 500)
	log.Error(err)
}
