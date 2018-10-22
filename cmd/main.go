package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/rustyeddy/moni"
	"github.com/rustyeddy/store"
)

var (
	config *moni.Configuration
)

func init() {
	flag.StringVar(&config.Output, "output", "stdout", "Were to send log output")
	flag.StringVar(&config.Level, "level", "warn", "Log level to set")
	flag.StringVar(&config.Format, "format", "json", "Format to print log files")

	flag.StringVar(&config.Addrport, "addr", ":8888", " an Daemon in the background")

	// What "cmd" or "mode" to run the command crawl, run cli or daemon
	flag.BoolVar(&config.Cli, "cli", false, "Run a command line client")
	flag.BoolVar(&config.Serve, "serve", false, "Run as a service")

	flag.StringVar(&config.ConfigFile, "cfg", "/srv/moni/config.json", "Use configuration file")

	flag.IntVar(&config.Depth, "depth", 1, "Max crawl depth")
	flag.StringVar(&config.Pubdir, "dir", "pub", "Serve the site from this dir")
	flag.StringVar(&config.StoreDir, "store", "/srv/moni/", "Directory for Store to use")

	flag.BoolVar(&config.Profile, "prof", false, "Profile our http server (daemon)")
}

var (

	// See config.go for all configuration variables and all the flags
	config  Configuration
	storage *store.Store
)

func getStorage() *store.Store {
	var err error
	if storage == nil {
		if storage, err = store.UseStore(config.StoreDir); err != nil {
			log.Fatalf("Fataling getting our storage %s", config.StoreDir)
		}
	}
	return storage
}

func getStoredObject(name string, obj interface{}) error {
	st := getStorage()
	_, err := st.FetchObject(name, obj)
	if err != nil {
		return errors.New("failed to find " + name)
	}
	return nil
}

func main() {
	var (
		// We will have a server if we are running in background,
		// we will have a client if we are running in the foreground.
		// If we are running in the foreground, we may also be running
		// in the background.  And you guessed it!  If we are running in
		// the background, we may also be running in the foreground.
		//
		// We may also be running a *one time* command in the foreground,
		// in this case we won't be running in the background.
		//
		// All commands, both REST (server) & Cli (foreground) take the
		// same trip path through the code.  The string from the Cli
		// command are basically concatenated with the appropriate
		// characters that form a valid REST url (endpoint and args).
		//
		// For example the following obviously crawls the site, we'll
		// simply transform the cmd line into the corresponding url.
		//
		//    % moni crawl example.com  ~~> /crawl/example.com
		//
		// With the url (end point) of a REST request, we'll create an
		// http.Request and an http.ResponseWriter allowing us to
		// pass our cmdline url to directly to the handler.!!.  Almost,
		// turns out, before we pass the request to the request handler
		// we need to do some argument handling, specifically get the
		// arguments from the request into the mux vars variable (mouth
		// full of marbles).
		//
		// No problem, http.
		//

		// With the router and http handler, only the cli calls the router
		// and handler directly avoiding the TCP RTT.
		srv *moni.Server
		cli *moni.Client
	)

	// Flags are all handled in the config.go file just because there
	// are lots of them with some post processing bits, perfer to keep
	// the flow of main clean, though a quick look at configs and flags
	// in config.go will be useful
	flag.Parse()

	// ====================================================================
	// We gotta put/persist data somewhere, eventually, right?!. The setup
	// will default to local storage (if we have it).  If we are running as
	// a serverless app, we may not have a local filesystem, hence s3, gcp
	// or DO spaces(?) must be configured
	st := getStorage()

	// ====================================================================
	// Read our configurations and various Data if they exist.  If they
	// do NOT exist, we will start with reasonable defaults.  If config
	// ReadFile() failes the program will fail.  Refuse to run with a
	// broken configuration.
	config.ReadFile()

	// ====================================================================
	// Create and run the server if the program is supposed to
	if config.Serve {
		srv = httpServer()
		go startServer(srv)
	}

	// ====================================================================
	// Create loop in a command prompt performing what ever is needed
	if config.Cli {
		cli = NewClient(os.Stdout, config.Addrport)
		go cli.Start()
	}

	// ========================================================================
	// Process commands from the command line
	nargs := len(flag.Args())
	if nargs > 0 {
		var resp *http.Response
		cmd := flag.Arg(0)

		// What address and port to connect to and something to
		// write the output to (io.Writer)
		cli := NewClient(os.Stdout, config.Addrport)
		arg := flag.Arg(1)

		// Run a single command in the foreground
		switch cmd {
		case "GET", "POST", "PUT", "DELETE":
			resp = cli.Do(cmd, arg)
			body := GetBody(resp)
			fmt.Fprintf(cli.Writer, "%s %s\n", cmd, arg)
			fmt.Fprintf(cli.Writer, "\t%s\n", string(body))
		case "crawl":
			cli.CrawlUrl(arg)
		case "crawlids":
			cli.CrawlList()
		case "view":
			cli.CrawlId(arg)
		}
	}

	// ====================================================================
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	if config.Cli || config.Serve {

		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)

		// Block until we receive our signal.
		<-c

		// Start cleaning up before we shut down.  Make sure we flush all data
		// we want flushed ...
		st.Shutdown()

		// Create a deadline to wait for.
		ctx, cancel := context.WithTimeout(context.Background(), config.Wait)
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
