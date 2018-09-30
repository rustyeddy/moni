package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

// HTTPServer will register routes, open connection and listen for incoming
func HTTPServer(addrport string) {

	r := mux.NewRouter()
	r.Handle("/", http.FileServer(http.Dir("./public")))
	r.HandleFunc("/info", handleInfo)
	r.HandleFunc("/walk/{url}", handleWalk)

	// listen for connections...
	log.Infoln("listening on ", addrport)
	log.Fatal(http.ListenAndServe(addrport, r))
}

// Send back diagnostic info
func handleInfo(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Handling an HTML request ")

	info := make(map[string]string)
	info["Routes"] = "/, /info, /walk/{url}"
	msg, err := json.Marshal(&info)
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	fmt.Fprint(w, msg)
}
