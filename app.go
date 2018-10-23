package moni

import (
	"net/http"

	"github.com/rustyeddy/moni/app"
)

// AppHandler will compose a response to the request
// and in the process most likely need to gather a few peices
// of information, put them together in the right order and
// send back to the caller
func AppHandler(w http.ResponseWriter, r *http.Request) {

	b := app.NewBuilder()

	data := struct {
		Title string
		Name  string

		// Maybe?
		Crawls []string
		Sites  Sitemap
		Pages  Pagemap
	}{
		Title:  "ClowdOps",
		Crawls: GetCrawls(),
		Sites:  GetSites(),
		Pages:  GetPages(),
	}

	b.Assemble(w, "index.html", data)
}
