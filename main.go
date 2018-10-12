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

	Pages Pagemap = make(Pagemap)
	Sites Sitemap = make(Sitemap)

	ACL AccessList = AccessList{
		Allowed:     make(map[string]int),
		Rejected:    make(map[string]int),
		Unsupported: make(map[string]int),
	}
)

func main() {
	flag.Parse()

	// Flags are all handled in the config.go file just because there
	// are lots of them with some post processing bits, perfer to keep
	// the flow of main clean, though a quick look at configs and flags
	// in config.go will be useful

	// ====================================================================
	// Setup storage (default local storage, look for redis or mongo)
	Storage := initStorage(Config.StoreDir)
	AssertNotNil(Storage)

	// ====================================================================
	// Create and run the server if the program is supposed to
	var srv *http.Server
	if Config.Serve {
		srv = httpServer()
		go startServer(srv)
	}

	// ====================================================================
	// Create loop in a command prompt performing what ever is needed
	var cli *Client
	if Config.Cli {
		cli = NewClient(Config.Addrport)
		go cli.Start()
	}

	// ========================================================================
	// Process commands from the command line
	if len(os.Args) > 0 {
		// Run a single command in the foreground
		switch os.Args[0] {
		case "crawl":
			cli := NewClient(Config.Addrport)
			go cli.CrawlUrls(os.Args)
		}
	}

	// ====================================================================
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	if Config.Cli || Config.Serve {

		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)

		// Block until we receive our signal.
		<-c

		// Start cleaning up before we shut down.  Make sure we flush all data
		// we want flushed ...
		Storage.Shutdown()

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

	}
	os.Exit(0)
}

func ShutdownHandler(w http.ResponseWriter, r *http.Request) {
	// Time to write out all open files
	os.Exit(2)
}

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

func startServer(srv *http.Server) (err error) {
	// Run our server in a goroutine so that it doesn't block.
	if err = srv.ListenAndServe(); err != nil {
		log.Println(err)
	}
	return err
}

func initStorage(dir string) (storage *store.Store) {
	// Setup Storage ~ depending on what we have configured we are
	// going to be reading and storing lots of stuff
	var err error
	if Storage, err = store.UseStore(Config.StoreDir); err != nil {
		log.Fatal("failed to use store", Config.StoreDir, err)
	}
	return storage
}
