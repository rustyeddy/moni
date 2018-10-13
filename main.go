package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

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

	var (
		srv *http.Server
		cli *Client
	)

	// Flags are all handled in the config.go file just because there
	// are lots of them with some post processing bits, perfer to keep
	// the flow of main clean, though a quick look at configs and flags
	// in config.go will be useful
	flag.Parse()

	// ====================================================================
	// Setup storage (default local storage, look for redis or mongo)
	Storage := initStorage(Config.StoreDir)
	AssertNotNil(Storage)

	// ====================================================================
	// Create and run the server if the program is supposed to
	if Config.Serve {
		srv = httpServer()
		go startServer(srv)
	}

	// ====================================================================
	// Create loop in a command prompt performing what ever is needed
	if Config.Cli {
		cli = NewClient(Config.Addrport)
		go cli.Start()
	}

	// ========================================================================
	// Process commands from the command line
	nargs := len(flag.Args())
	if nargs > 0 {
		fmt.Println("Crawling ...  ")

		// Run a single command in the foreground
		switch flag.Arg(0) {
		case "crawl":
			cli := NewClient(Config.Addrport)
			cli.CrawlUrls(os.Args)
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

func initStorage(dir string) (storage *store.Store) {
	// Setup Storage ~ depending on what we have configured we are
	// going to be reading and storing lots of stuff
	var err error
	if Storage, err = store.UseStore(Config.StoreDir); err != nil {
		log.Fatal("failed to use store", Config.StoreDir, err)
	}
	return storage
}
