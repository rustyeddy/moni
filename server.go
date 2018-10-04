package main

import (
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type Server struct {
	Addrport string // config option
	*mux.Router
	*log.Logger
}

// StartServer registers the handlers for the API routes
// and to return files from the static website.
func NewServer(addrport string) (srv *Server) {
	if addrport == "" {
		log.Fatalf("Server must have a port to listen with")
	}

	srv = &Server{
		Addrport: addrport,
		Logger:   log.New(),
	}

	log.Infoln("Creating HTTP Router")
	srv.Router = mux.NewRouter()
	srv.HandleFunc("/crawl/{url}", HandleCrawl)

	return srv
}

// Start your engines
func (srv *Server) Start(done chan<- bool) {
	log.Infoln("Server listening on ", srv.Addrport)

	// r := mux.NewRouter()
	// r.HandleFunc("/crawl/{url}", HandleCrawl)
	// r.HandleFunc("/crawl/{url}", HandleCrawl)
	http.Handle("/", srv.Router)

	//err := http.ListenAndServe(srv.Addrport, r)
	http.ListenAndServe(srv.Addrport, srv.Router)
	done <- true
}
