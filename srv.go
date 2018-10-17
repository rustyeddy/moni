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

	// Register the application url and handlers
	registerApp(r)

	// Regsiter the site crawler routers
	registerCrawlers(r)

	// Register the update handler
	registerUpdate(r)

	// Register the profiler
	registerProfiler(r) // make these plugins ...

	return createServer(r, Config.Addrport)
}

// registerCrawlers will register all handlers related to the
// crawling activities
func registerCrawlers(r *mux.Router) {
	r.HandleFunc("/acl", ACLHandler)               // Display ACLs
	r.HandleFunc("/crawlids", CrawlListHandler)    // Display "recent" crawl jobs
	r.HandleFunc("/crawl/{url}", CrawlHandler)     // Create a (recurring) crawl job for url
	r.HandleFunc("/crawlid/{cid}", CrawlIdHandler) // Display a specific crawl job
}

func registerUpdate(r *mux.Router) {
	r.HandleFunc("/update/", UpdateHandler) // Handle updates when signaled from github
}

// registerApp will register a static file handler allowing us to serve
// up the web pages for our application.
func registerApp(r *mux.Router) {
	// This will serve files under http://localhost:8000/static/<filename>
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(Config.Pubdir))))
	r.HandleFunc("/", AppHandler)
}

// registerProfiler makes the profiler available at the specified locations
func registerProfiler(r *mux.Router) {
	r.HandleFunc("/pprof/", pprof.Index)
	r.HandleFunc("/pprof/cmdline", pprof.Cmdline)
	r.HandleFunc("/pprof/profile", pprof.Profile)
	r.HandleFunc("/pprof/symbol", pprof.Symbol)
	r.HandleFunc("/pprof/trace", pprof.Trace)
}

// Create the Server setting the address, router and some timeouts
func createServer(r *mux.Router, addrport string) *http.Server {
	srv := &http.Server{
		Addr: addrport,
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
