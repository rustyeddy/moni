package moni

import (
	"fmt"
	"net/http"
)

// AppHandler will compose a response to the request
// and in the process most likely need to gather a few peices
// of information, put them together in the right order and
// send back to the caller
func AppHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, "request: %+v\n", r)
	/*

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
	*/
}
