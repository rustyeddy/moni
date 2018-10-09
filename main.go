package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/rustyeddy/inv/store"
)

var (

	// See config.go for specific configuration variables
	Config  Configuration
	Storage *store.Store

	Pages PageInfomap = make(PageInfomap)
	Sites Sitemap     = make(Sitemap)
	ACL   AccessList
)

func init() {
	ACL = make(map[string]bool)
}

func main() {
	flag.Parse()

	r := mux.NewRouter()
	r.HandleFunc("/", AppHandler)
	r.HandleFunc("/crawl/{url}", CrawlHandler)
	r.HandleFunc("/update/", UpdateHandler)

	// Set the profile handlers if we have flagged them to be turned on
	if Config.Profile {
		r.HandleFunc("/debug/pprof/", pprof.Index)
		r.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		r.HandleFunc("/debug/pprof/profile", pprof.Profile)
		r.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
		r.HandleFunc("/debug/pprof/trace", pprof.Trace)
	}

	var err error
	if Storage, err = store.UseStore(Config.StoreDir); err != nil {
		log.Fatalf("failed to use store", Config.StoreDir, err)
	}

	srv := &http.Server{
		Addr: "0.0.0.0:8888",
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r, // Pass our instance of gorilla/mux in.
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), Config.Wait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("shutting down")
	os.Exit(0)
}

func ShutdownHandler(w http.ResponseWriter, r *http.Request) {
	// Time to write out all open files
	os.Exit(2)
}
