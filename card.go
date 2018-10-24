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

func (c *Card) WriteHTML(w http.ResponseWriter, tmpl *template.Template) {
	err := tmpl.ExecuteTemplate(w, "card.html", c)
	if err != nil {
		log.Fatalln(err)
	}
}
