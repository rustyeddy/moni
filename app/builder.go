package app

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

type PageData struct {
	Title string // name of the page (url title)
	Name  string // name for fun and profit
	Tmpl  string // name
	Frag  string // request.URL.Fragment

	Rows  []*Row
	Cards []*Card
}

// Builder constructs (and sends) the response back to the
// user.  It determines with template pieces to put together,
// assembles them and off they go
type PageBuilder struct {
	TemplateBasedir string
	TemplateName    string
	*template.Template

	*PageData
}

// NewBuilder will find and compile the templates, which are broke
// into layout (comprise the structure of the site) and partials
// (comprise content elements)
func NewPageBuilder() *PageBuilder {
	pb := PageBuilder{
		TemplateBasedir: "../app/tmpl",
		TemplateName:    "index.html", // not actually used .?.
		PageData:        &PageData{Title: "Clowd ~ Ops "},
	}
	pb.PageData.Name = "Rusty Eddy"
	pb.PrepareTemplates()
	return &pb
}

func (pb *PageBuilder) PrepareTemplates() {
	pattern := filepath.Join(pb.TemplateBasedir, "*.html")
	pb.Template = template.Must(template.ParseGlob(pattern))
	pb.DumpTemplates()
}

func (b *PageBuilder) DumpTemplates() {
	fmt.Println("Templates: ", b.Template.Name())
	fmt.Println(b.DefinedTemplates())
}

func (b *PageBuilder) AddCard(w http.ResponseWriter, card *Card) {
	b.Cards = append(b.Cards, card)
}

// Assemble traverses our local representation of the outgoing documents,
// occaisionally run stuff through a template, writing out successful
// stuff as required.
func (b *PageBuilder) Assemble(w http.ResponseWriter, tmplname string) {

	// Here we go, create our html for our site.  Building the page happens
	// in two parts.  1. A semi-generic frame is created with designated areas
	// can be overwritten with application specific information.
	if err := b.ExecuteTemplate(w, "index.html", b); err != nil {
		log.Fatalln(err)
	}
}
