package moni

import (
	"net/http"

	"github.com/rustyeddy/moni/app"
	log "github.com/sirupsen/logrus"
)

type AppData struct {
	Name    string
	Frag    string
	Content string
	Title   string

	Crawls  []*moni.Crawls
	Sites   *moni.Sitemap
	Pages   *moni.Pages
	Domains *moni.Domains
}

// AppHandler will compose a response to the request
// and in the process most likely need to gather a few peices
// of information, put them together in the right order and
// send back to the caller
func AppHandler(w http.ResponseWriter, r *http.Request) {

	log.Infoln("AppHandler request: %+v", r)

	b := app.NewBuilder()
	d := AppData{
		Name:    "Rusty Eddy",
		Frag:    r.URL.Fragment,
		Content: "Hello, World!",
		Title:   "Clowd Ops",

		Crawls:  moni.GetCrawls(),
		Sites:   monilGetSites(),
		Pages:   moni.GetPages(),
		Domains: moni.GetDomains(),
	}
	b.Assemble(w, "index.html", data)
}
