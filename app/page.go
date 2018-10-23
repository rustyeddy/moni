package app

import (
	"fmt"
	"html/template"
	"net/http"
)

type PageData struct {
	Title string // name of the page (url title)
	Tmpl  string // name
}

// Builder constructs (and sends) the response back to the
// user.  It determines with template pieces to put together,
// assembles them and off they go
type PageBuilder struct {
	Layouts *template.Template
	//Partials *template.Template
}

// NewBuilder will find and compile the templates, which are broke
// into layout (comprise the structure of the site) and partials
// (comprise content elements)
func NewBuilder() (b *PageBuilder) {
	b = new(PageBuilder)
	b.PrepareTemplates()
	return b
}

// NewBuilder will find and compile the templates, which are broke
// into layout (comprise the structure of the site) and partials
// (comprise content elements)
func (b *PageBuilder) PrepareTemplates() {
	layouts := template.Must(template.ParseGlob("../app/tmpl/*.html"))
	b.Layouts = layouts
	b.DumpTemplates()
}

func (b *PageBuilder) DumpTemplates() {
	fmt.Println("layouts: ", b.Layouts.Name())
	fmt.Println(b.Layouts.DefinedTemplates())
}

// Assemble the template with data provide
func (b *PageBuilder) Assemble(w http.ResponseWriter, name string, data interface{}) {
	err := b.Layouts.ExecuteTemplate(w, "index.html", data)
	if err != nil {
		fmt.Fprintf(w, "internal error -> %+v", err)
		return
	}
}
