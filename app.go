package main

import (
	"html/template"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func AppHandler(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Title  string
		Body   string
		Header string
	}{"Willy Wonkers", "My Sexy Body", "Head Case Industries"}

	var t = template.Must(template.ParseGlob("tmpl/*.html"))
	if err := t.Execute(w, data); err != nil {
		log.Errorf("Template failed %v", err)
		JSONError(w, err)
	}
}
