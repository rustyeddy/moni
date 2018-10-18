package main

import (
	"fmt"
	"html/template"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func AppHandler(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Title  string
		Crawls []string
	}{"Willy Wonkers", nil}

	base := "dash/tmpl/"
	tmplts := []string{
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

	log.Infoln("Request received for AppHandler")

	data.Crawls = GetCrawls()
	var t = template.Must(template.ParseFiles(tmplts...))
	if err := t.Execute(w, data); err != nil {
		log.Errorf("PUKE Template failed %v", err)
		fmt.Fprintln(w, "interal error")
	}
}
