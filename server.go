package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

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

	srv.HandleFunc("/crawl/{url}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		url := vars["url"]

		if !strings.HasPrefix("http", url) {
			url = "http://" + url
		}

		srv.Infoln("crawl", r.URL)
		page, err := Crawl(url)
		if err != nil {
			fmt.Fprintf(w, "url", err)
		}

		jbytes, err := json.Marshal(page)
		if err != nil {
			srv.Errorln("marshal json", url, err)
		}

		w.Header().Add("Content-Type", "application/json")
		w.Write(jbytes)
	})

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
