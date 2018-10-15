package main

import (
	"net/http"
	"net/http/pprof"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

// httpServer creates the router, registers the handlers then
// creates the server, primed and ready to be started.  We pass the
// *http.Server back to the caller, allowing it (main as of this
// writing) it to `go startServer(srv)` start the server as a
// Go Routine().
func httpServer() *http.Server {
	r := mux.NewRouter()
	r.HandleFunc("/", AppHandler)
	r.HandleFunc("/crawl/{url}", CrawlHandler)
	r.HandleFunc("/acl", ACLHandler)
	r.HandleFunc("/update/", UpdateHandler)

	// Set the profile handlers if we have flagged them to be turned on
	if Config.Profile {
		r.HandleFunc("/debug/pprof/", pprof.Index)
		r.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		r.HandleFunc("/debug/pprof/profile", pprof.Profile)
		r.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
		r.HandleFunc("/debug/pprof/trace", pprof.Trace)
	}

	// Create the Server itself
	srv := &http.Server{
		Addr: Config.Addrport,
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r, // Pass our instance of gorilla/mux in.
	}
	return srv
}

// startServer actually starts the server, most likely as a Go routine
// background function.
func startServer(srv *http.Server) (err error) {
	// Run our server in a goroutine so that it doesn't block.
	if err = srv.ListenAndServe(); err != nil {
		log.Println(err)
	}
	return err
}
