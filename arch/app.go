package moni

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"

	"github.com/gorilla/mux"
	"github.com/rustyeddy/store"
	log "github.com/sirupsen/logrus"
)

// All global variables
var (
	app    *App
	server *http.Server
	router *mux.Router
	sites  Sitemap
	pages  Pagemap
)

// ====================================================================
//                           App
// ====================================================================

// App is the One TRUE App! All Hail the App!  It is the global context
// of everything.  It contains some information for the web page to be
// displayed, it also maintains configurations and managers the server
// and the scheduler.
type App struct {
	AddrPort string

	// Basic meta stuff for App web page and content
	Title string // name of the page (url title)
	Name  string // name for fun and profit
	Tmpl  string // base or frame template
	Frag  string // request.URL.Fragment

	// This is only real configuration
	Configuration

	// Tmplates to handle html and text formatting
	AppTemplates
	*log.Entry
}

// NewApp will produce a new App
func GetApp(cfg *Configuration) (app *App) {
	if app == nil {
		app = &App{
			Name:  "ClowOpsApp",
			Title: "Clowd ~ Operations",
			Tmpl:  "index.html",
		}
		app.Title = app.Name
		app.Configuration = *cfg

		// Setup the logger
		app.Entry = log.WithFields(log.Fields{
			"app":  app.Title,
			"tmpl": app.Tmpl,
		})
	}
	// TODO: setup router with subrouting
	// code here xxx
	return app
}

func (app *App) Init() *App {
	// Some global could make them the apps

	st, err := store.UseFileStore("./etc")
	IfErrorFatal(err)

	acl = initACL()
	sites = initSites()
	pages = initPages()

	// urlQ = NewURLQ()
	// crawlQ = NewCrawlQ()
	// saveQ = NewSaveQ()

	// Create the server ~ And register the app
	server, router = GetServer(app.AddrPort)
	app.Register(router)
	return app
}

func (app *App) StartService() {
	//go urlQ.Watch()
	//go crawlQ.Watch()
	//go saveQ.Watch()
	server.ListenAndServe()
}

// ====================================================================
//                      App Templates
// ====================================================================

// Contains various pointers to Go templates and the compiled
// version of the templates.
type AppTemplates struct {
	TmplBasedir string
	TmplName    string
	*template.Template
}

// Acculmulate the data needed for the template
type AppData struct {
	Sites []*Site
	*Configuration
}

// Builder constructs (and sends) the response back to the
// user.  It determines with template pieces to put together,
// assembles them and off they go
func (app *App) PrepareTemplates(tmpldir string) {
	pattern := filepath.Join(tmpldir, "*.html")
	app.Template = template.Must(template.ParseGlob(pattern))
}

func (app *App) DumpTemplates() {
	fmt.Println("Templates: ", app.Template.Name())
	fmt.Println(app.DefinedTemplates())
}

// ====================================================================
//                      App Assembler
// ====================================================================

// Assemble traverses our local representation of the outgoing documents,
// occaisionally run stuff through a template, writing out successful
// stuff as required.
func (app *App) Assemble(w http.ResponseWriter, tmplname string) {
	// Here we go, create our html for our site.  Building the page happens
	// in two parts.  1. A semi-generic frame is created with designated areas
	// can be overwritten with application specific information.
	if app.Template == nil {
		app.PrepareTemplates(app.Tmpldir)
	}
	d := &AppData{
		Configuration: &app.Configuration,
	}
	if err := app.ExecuteTemplate(w, "index.html", d); err != nil {
		app.Fatalln(err)
	}
}

func (app *App) Run() {
	// Wait for the server (and/or) client to end
	// ====================================================================
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), app.Wait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	app.Shutdown(ctx)

	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("shutting down")

	os.Exit(0)
}
