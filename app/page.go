package app

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
)

type PageData struct {
	Title string // name of the page (url title)
	Tmpl  string // name
}

// Builder constructs (and sends) the response back to the
// user.  It determines with template pieces to put together,
// assembles them and off they go
type PageBuilder struct {
	TemplateBasedir string
	TemplateName    string
	*template.Template
}

// NewBuilder will find and compile the templates, which are broke
// into layout (comprise the structure of the site) and partials
// (comprise content elements)
func NewPageBuilder() *PageBuilder {
	pb := PageBuilder{
		TemplateBasedir: "../app/tmpl",
		TemplateName:    "index.html",
	}
	pb.PrepareTemplates()
	return &pb
}

func (pb *PageBuilder) PrepareTemplates() {
	pattern := filepath.Join(pb.TemplateBasedir, "*.html")
	pb.Template = template.Must(template.ParseGlob(pattern))
}

func (b *PageBuilder) DumpTemplates() {
	fmt.Println("Templates: ", b.Name())
	fmt.Println(b.DefinedTemplates())
}

// Assemble the template with data provide
func (b *PageBuilder) Assemble(w http.ResponseWriter, name string, data *AppData) {

	err := b.ExecuteTemplate(w, name, data)
	if err != nil {
		fmt.Fprintf(w, "internal error -> %+v", err)
		return
	}
}
