package main

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"

	log "github.com/sirupsen/logrus"
)

var (
	tmplpat      string // Make this a flag ?
	tmplCompiled *template.Template
)

func init() {
	tmplpat = "app/tmpl/*.html"
}

func getCompiledTemplates() (t *template.Template) {
	if tmplCompiled == nil {
		tmplCompiled = template.Must(template.ParseGlob(tmplpat))
	}
	return tmplCompiled
}

func AppHandler(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Title  string
		Crawls []string
		Sites  Sitemap
		Pages  Pagemap
	}{
		Title:  "ClowdOps",
		Crawls: GetCrawls(),
		Sites:  GetSites(),
		Pages:  GetPages(),
	}

	var t *template.Template
	if t = getCompiledTemplates(); t == nil {
		JSONError(w, errors.New("failed to compile templates"))
		return
	}

	if err := t.Execute(w, data); err != nil {
		log.Errorf("Template PUKED %v ", err)
		fmt.Fprintf(w, "template BARFED %v", err)
	}
}
