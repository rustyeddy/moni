package moni

import (
	"html/template"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// CardOptions
type Cardmap map[string]string

func NewCardmap() (cm Cardmap) {
	cm = Cardmap{
		"title":    "",
		"subtitle": "",
		"image":    "",
		"cols":     "",
		"links":    "",
		"button":   "",
		"text":     "",
	}
	return cm
}

// Card represents a card from our Bootstrap Dashboard
type Card struct {
	Image    string   // Optional header image
	Title    string   // Title recomended
	Subtitle string   // Optional
	Cols     string   // Used to set the number of BootStrap colums
	Links    []string // array of links, if we have any
	Button   string   // optional button if we need one
	Text     string
}

/*
func (c *Card) Image() string {
	if c.image != "" {
		return `<img class='image' src='` + c.image + `' />`
	}
	return ""
}

func (c *Card) Header() string {
	if c.title != "" {
		return `<div class='card-header'>` +
			`<h5 class='card-title mb-0'>` + c.title + `</h5>` +
			`</div> <!-- card-header --> `
	}
	return ""
}

func (c *Card) Subtitle() string {
	if c.title != "" {
		return `<div class='card-subtitle text-muted'>` +
			`<h6 class='card-title mb-0'>` + c.title + `</h6>` +
			`</div> <!-- card-subtitle --> `
	}
	return ""
}

func (c *Card) Links() string {
	if c.links == nil {
		return ""
	}
	html := ""
	for _, l := range c.links {
		html = html + `<a class='card-link' href='` + l + `'>` + l + `</a> `
	}
	return html
}

func (c *Card) Button() string {
	if c.button == "" {
		return ""
	}
	return `<a href="{{ .Link }}" class='btn btn-primary'>{{ .Text }}</a>`
}
*/

func (c *Card) WriteHTML(w http.ResponseWriter, tmpl *template.Template) {
	err := tmpl.ExecuteTemplate(w, "card.html", c)
	if err != nil {
		log.Fatalln(err)
	}
}
