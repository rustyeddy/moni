package main

import (
	"html/template"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func AppHandler(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Title  string
		Crawls []string
	}{"Willy Wonkers", nil}

	data.Crawls = GetCrawls()
	var t = template.Must(template.ParseGlob("dash/tmpl/*.html"))
	if err := t.Execute(w, data); err != nil {
		log.Errorf("PUKE Template failed %v", err)
		JSONError(w, err)
	}
}
