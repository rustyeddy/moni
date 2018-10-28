package moni

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

var (
	app *App
)

func init() {
	app = NewApp(&DefaultConfig)
}

// ====================================================================
//                           App
// ====================================================================

// App is the One TRUE App! All Hail the App!  It is the global context
// of everything.  It contains some information for the web page to be
// displayed, it also maintains configurations and managers the server
// and the scheduler.
type App struct {

	// Basic meta stuff for App web page and content
	Title string // name of the page (url title)
	Name  string // name for fun and profit
	Tmpl  string // base or frame template
	Frag  string // request.URL.Fragment

	// Tmplates to handle html and text formatting
	AppTemplates

	*log.Entry
}

// NewApp will produce a new App
func NewApp(cfg *Configuration) (app *App) {
	app = &App{
		Name:  "ClowOpsApp",
		Title: "Clowd ~ Operations",
		Tmpl:  "index.html",
	}
	app.Title = app.Name
	SetConfig(cfg)

	// Setup the logger
	app.Entry = log.WithFields(log.Fields{
		"app":  app.Title,
		"tmpl": app.Tmpl,
	})
	return app
}

// NewApp will produce a new App
func NewTestApp(config *Configuration) (app *App) {
	app = &App{
		Name:  "ClowOpsApp",
		Title: "Clowd ~ Operations",
		Tmpl:  "index.html",
	}
	app.Title = app.Name
	return app
}

func (app *App) Start() {
	//StartDatabase()
	StartServer()
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
type Appdata struct {
	*Configuration
}

// Builder constructs (and sends) the response back to the
// user.  It determines with template pieces to put together,
// assembles them and off they go
func (app *App) PrepareTemplates(tmpldir string) {
	pattern := filepath.Join(tmpldir, "*.html")

	fmt.Printf("Reading templates dir %s from %s\n", tmpldir, pattern)
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
		app.PrepareTemplates(config.Tmpldir)
	}

	d := &Appdata{
		Configuration: config,
	}

	if err := app.ExecuteTemplate(w, "index.html", d); err != nil {
		app.Fatalln(err)
	}
}
