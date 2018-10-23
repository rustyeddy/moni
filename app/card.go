package app

import (
	"html/template"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// CardOptions
type CardOptions map[string]string

// Card represents a card from our Bootstrap Dashboard
// CSS Template
type Card struct {
	*Image
	Title string  // header
	Text  string  // body
	Links []*Link //

	TmplName string // Template name we'll use to convert to html
}

type Link struct {
	Text string
	Link string
}

func linksFromString(str string) (l []*Link) {
	return nil
}

// NewCard creates a new card
func NewCard(config map[string]string) (c *Card) {
	c = new(Card)

	for k, v := range config {
		switch k {
		case "Title":
			c.Title = v
		case "Text":
			c.Text = v
		case "Links":
			c.Links = linksFromString(v)
		case "Image":
			c.Image.Src = v
		case "Alt":
			c.Image.Alt = v
		default:
			log.Errorf("unexpected card part %s, ignoring .. %s with value ", k, v)
		}
	}
	return c
}

func (c *Card) WriteHTML(w http.ResponseWriter, tmpl *template.Template) {
	err := tmpl.ExecuteTemplate(w, "card.html", c)
	if err != nil {
		log.Fatalln(err)
	}
}
