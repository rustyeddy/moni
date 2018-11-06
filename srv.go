package moni

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/http/pprof"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

// ====================================================================
//                      Start Server
// ====================================================================

func (app *App) Shutdown(ctx context.Context) {
	app.Shutdown(ctx)
}

// httpServer creates the router, registers the handlers then
// creates the server, primed and ready to be started.  We pass the
// *http.Server back to the caller, allowing it (main as of this
// writing) it to `go startServer(srv)` start the server as a
// Go Routine().
func NewServer(addrport string) (s *http.Server, r *mux.Router) {
	r = mux.NewRouter()

	// The ACL handler
	registerACLHandler(r)

	// Register the site handler
	registerSites(r)

	// Register the pages
	registerPages(r)

	// Regsiter the site crawler routers
	registerCrawlers(r)

	// Register the update handler
	registerUpdate(r)

	// Just register the profiler
	registerProfiler(r) // make these plugins ...

	// Create the Server setting the address, router and some timeouts
	s = &http.Server{
		Addr: addrport,
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r, // Pass our instance of gorilla/mux in.
	}
	return s, r
}

// RequestLoopback will pass the given handler and url to a special
// http test client and and "server", bypassing the network.  The
// Request, URL and variables are sanitized and prepared before
// calling the Handler function (just like a real server).  The request
// returned by the Handler is returned to the calling test for
// examination.  The calling test has the context required to determine
// the passability of the test.
func ServiceLoopback(h http.HandlerFunc, verb string, url string) *http.Response {

	// Craft up a request with the URL we want to test
	req := httptest.NewRequest(verb, url, nil)
	w := httptest.NewRecorder()

	// Do not give the writer and request to the handler directly.  The args
	// will not have been processed.  Register it as a handler the let mux
	// parse our the args and setup other important things, then let it
	// call the handler itself.
	_, r := NewServer(":8888")
	//r := srv.Handler

	// This will cause the actual crawling, CrawlHandler will be called
	// with all the appropriate header and argument processing.
	r.ServeHTTP(w, req)

	// Look at the response (which had been sent to the client).
	resp := w.Result()
	if resp == nil {
		log.Errorln("failed to get a response")
	}
	return resp
}

func registerUpdate(r *mux.Router) {
	r.HandleFunc("/update/", UpdateHandler) // Handle updates when signaled from github
}

// registerProfiler makes the profiler available at the specified locations
func registerProfiler(r *mux.Router) {
	r.HandleFunc("/pprof/", pprof.Index)
	r.HandleFunc("/pprof/cmdline", pprof.Cmdline)
	r.HandleFunc("/pprof/profile", pprof.Profile)
	r.HandleFunc("/pprof/symbol", pprof.Symbol)
	r.HandleFunc("/pprof/trace", pprof.Trace)
}
