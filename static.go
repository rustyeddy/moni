package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type StaticServer struct {
	Addrport string
	Basedir  string

	*mux.Router
	*log.Logger
}

func NewStaticServer(addrport string) (srv *StaticServer) {
	if addrport == "" {
		log.Fatalf("Server must have a port to listen with")
	}
	srv = &StaticServer{
		Addrport: addrport,
		Logger:   log.New(),
		Basedir:  Config.Pubdir,
	}
	srv.Router = mux.NewRouter()
	srv.Handle("/", http.FileServer(http.Dir("./pub")))
	return srv
}

func (srv *StaticServer) Start(done chan<- bool) {
	log.Infoln("Server listening on ", srv.Addrport)
	err := http.ListenAndServe(srv.Addrport, srv.Router)
	if err != nil {
		log.Errorf("Server terminated error %v", err)
	}
	done <- true
}

func HandleFile(w http.ResponseWriter, fname string) {
	fmt.Fprintf(w, fname)
}
