package moni

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

// RegisterApp will register a static file handler allowing us to serve
// up the web pages for our application.
func registerApp(r *mux.Router) {
	// This will serve files under http://localhost:8888/static/<filename>
	staticdir := "../static"
	if _, err := os.Stat(staticdir); os.IsNotExist(err) {
		staticdir = "static"
		if _, err := os.Stat(staticdir); os.IsNotExist(err) {
			IfErrorFatal(err, "registerApp")
		}
	}
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(staticdir))))
	r.HandleFunc("/", AppHandler)
}

// AppHandler will compose a response to the request
// and in the process most likely need to gather a few peices
// of information, put them together in the right order and
// send back to the caller
func AppHandler(w http.ResponseWriter, r *http.Request) {

	app.Title = "Application Interface"
	app.Name = "Rusty Eddy"
	app.Frag = r.URL.Fragment

	app.Assemble(w, "index.html")
}

// AppHandler will compose a response to the request
// and in the process most likely need to gather a few peices
// of information, put them together in the right order and
// send back to the caller
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	app.Title = "Home Handler"
	app.Name = "Rusty Eddy"
	app.Frag = r.URL.Fragment

	app.Assemble(w, "home.html")
}
