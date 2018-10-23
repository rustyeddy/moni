package moni

import (
	"net/http"

	"github.com/rustyeddy/moni/app"
	log "github.com/sirupsen/logrus"
)

// AppHandler will compose a response to the request
// and in the process most likely need to gather a few peices
// of information, put them together in the right order and
// send back to the caller
func AppHandler(w http.ResponseWriter, r *http.Request) {
	log.Infof("AppHandler request: %+v", r)

	pb := app.NewPageBuilder()
	d := app.AppData{
		Title:   "Application Interface",
		Name:    "Rusty Eddy",
		Frag:    r.URL.Fragment,
		Content: "<h2>Content</h2>",
	}
	pb.Assemble(w, pb.TemplateName, &d)
}
