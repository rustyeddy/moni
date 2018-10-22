package main

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"

	log "github.com/sirupsen/logrus"
)

var (
	compiledTemplates *template.Template
	templateFiles     []string
)

func init() {
	base := "dash/tmpl/"
	templateFiles = []string{
		base + "index.html",
		base + "head-cheese.html",
		base + "sidebar.html",
		base + "header-nav.html",
		base + "nav-messages.html",
		base + "nav-alerts.html",
		base + "nav-profile.html",
		base + "recent-crawls.html",
		base + "crawl-details.html",
		base + "site-list.html",
	}
}

func getCompiledTemplates() (t *template.Template) {
	if compiledTemplates == nil {
		compiledTemplates = template.Must(template.ParseFiles(templateFiles...))
	}
	return compiledTemplates
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
