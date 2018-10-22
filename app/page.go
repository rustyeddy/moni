package app

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
)

type Page struct {
	Title string // name of the page (url title)
	Tmpl  string // name
	Data  map[string]string
}

// Builder constructs (and sends) the response back to the
// user.  It determines with template pieces to put together,
// assembles them and off they go
type PageBuilder struct {
	*template.Template
}

// NewBuilder will find and compile the templates, which are broke
// into layout (comprise the structure of the site) and partials
// (comprise content elements)
func NewBuilder() (b *PageBuilder) {
	b = new(PageBuilder)

	// Following is from the Go code ...
	// We load all the templates before execution. This package does not require
	// that behavior but html/template's escaping does, so it's a good habit.
	for _, d := range []string{"app/tmpl", "app/tmpl/partials"} {
		b.PrepareTemplates(d)
	}
	return b
}

// Assemble the template with data provide
func (b *PageBuilder) Assemble(w http.ResponseWriter, name string, data interface{}) {
	// err := b.TemplateExectuteTemplate(w, name, data, nil)
	err := b.Template.Execute(w, data)
	if err != nil {
		fmt.Fprintf(w, "internal error\n")
		return
	}
}

// NewBuilder will find and compile the templates, which are broke
// into layout (comprise the structure of the site) and partials
// (comprise content elements)
func (b *PageBuilder) PrepareTemplates(tmpldir string) {
	// We are going to first, find and compile all the .html "layout"
	// files (templates) in config.Tmpldir.  We will then create
	// a new set of templates from the "partials directory".

	// pattern is the glob pattern used to find all the template files.
	pattern := filepath.Join(tmpldir, "*.html")
	b.Template = template.Must(template.ParseGlob(pattern))
}
