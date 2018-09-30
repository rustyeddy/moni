package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type pageInfo struct {
	StatusCode int
	Links      map[string]int
}

func Server(addrport string) {
	r := mux.NewRouter()
	r.HandleFunc("/", handleHome)
	r.HandleWalk("/walk/", handleWalk)

	log.Infoln("listening on ", addrport)
	log.Fatal(http.ListenAndServe(addrport, nil))
}

// Send back the home page
func handleHome(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Handling an HTML request ")
}
