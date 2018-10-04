package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	// Parse command line args setting config values
	// as set in config.go
	flag.Parse()
	start2()
}

func start2() {
	r := mux.NewRouter()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Wee I go!")
	})
	r.HandleFunc("/crawl/{url}", HandleCrawl)
	http.Handle("/", r)
	http.ListenAndServe(":4444", r)
}

func start1() {
	// Declare the done channel to communicate when the
	// server has completed
	var chHttp, chStatic, chClient chan bool

	// Figure what command we are going to run
	// with what specific arguments.
	if Config.Client {
		cli := NewClient(Config.HttpAddr)
		chHttp = make(chan bool)
		go cli.Start(chHttp)
	}

	if Config.StartStatic {

		// TODO: make this a configuration item
		srv := NewStaticServer(Config.StaticAddr)
		chStatic = make(chan bool)
		go srv.Start(chStatic)
	}

	if Config.StartHttp {
		srv := NewServer(Config.HttpAddr)
		chHttp = make(chan bool)
		srv.Start(chHttp)
	}

	for {
		select {
		case _ = <-chHttp:
			close(chHttp)
		case _ = <-chStatic:
			close(chStatic)
		case _ = <-chClient:
			close(chClient)
			//case <-time.After(4 * time.Second):
		}
	}
	fmt.Println("Server has finished")
}
