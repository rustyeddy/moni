package moni

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

/*
 * The One TRUE App
 */
type App struct {

	// Basic meta stuff
	Title string // name of the page (url title)
	Name  string // name for fun and profit
	Tmpl  string // base or frame template
	Frag  string // request.URL.Fragment

	Configuration

	// Tmplates to handle html and text formatting
	AppTemplates

	// The following
	cards []*Card

	// We will Explicitly list our cards here
	*SitesCard
	*StorageCard
	LogCard *Card // Just a Card
}

type AppTemplates struct {
	TmplBasedir string
	TmplName    string
	*template.Template
}

// NewApp will produce a new App
func NewApp(name string) (app *App) {
	app = &App{
		Name:  name,
		Tmpl:  "index.html",
		Title: name, // default name, but can change
	}
	app.TmplBasedir = "../tmpl"
	app.TmplName = "index.html"
	app.Title = "Clowd ~ Operations"
	app.Name = "Rusty Eddy"

	// Prepare templates now
	app.PrepareTemplates()

	return app
}

// Builder constructs (and sends) the response back to the
// user.  It determines with template pieces to put together,
// assembles them and off they go
func (app *App) PrepareTemplates() {
	pattern := filepath.Join(app.TmplBasedir, "*.html")
	app.Template = template.Must(template.ParseGlob(pattern))
	if app.Debug {
		app.DumpTemplates()
	}
}

func (app *App) DumpTemplates() {
	fmt.Println("Templates: ", app.Template.Name())
	fmt.Println(app.DefinedTemplates())
}

// AddCard to the applicaiton
func (a *App) AddCard(c *Card) *App {
	a.cards = append(a.cards, c)
	return a // allow chaining
}

// Assemble traverses our local representation of the outgoing documents,
// occaisionally run stuff through a template, writing out successful
// stuff as required.
func (a *App) Assemble(w http.ResponseWriter, tmplname string) {
	// Here we go, create our html for our site.  Building the page happens
	// in two parts.  1. A semi-generic frame is created with designated areas
	// can be overwritten with application specific information.
	if err := a.ExecuteTemplate(w, "index.html", a); err != nil {
		log.Fatalln(err)
	}
}
