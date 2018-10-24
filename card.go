package moni

import (
	"html/template"
	"net/http"

	log "github.com/sirupsen/logrus"
)

/*

This Dashboard holds lots of different types of information, including
text, images, graphs, tables and lots of data.  The *Card* is the basic
container for all these types of information.

A quick look at the dashboard will make the cards obvious.  They present
a white background on the light grey Dashboard panel (background).

Cards have many optional elements including: title, subtitle, image, buttons,
links and maybe a footer.

Card content can vary from card to card, content maybe plain text, table,
lists, graphs, basically anything you can put into an HTML document.

Basically we will use the template system to create the framework of the
dashboard including the sidebar and header navs.  The main content area
will be a light grey background with a "Grid of cards".

We will create a card (or more) for every type of visual we want to see
and feed them to the page builder accordingly.

All cards need 3 things to produce a peice of visual content:

1. A HTML template defining the card layout.  Templates may have one or
   more "variables" defined to provide data in specific places

2. A data set that is used to replace variables with real data.  If
   the template requires variables/values that do not exist in a dataset
   the template will generate a *Runtime error*.

3. An http.ResponseWriter to write the resulting data to.

The first two can be prepared and compiled before hand.
*/

// Card represents a card from our Bootstrap Dashboard
type Card struct {
	Name     string
	Image    string   // Optional header image
	Title    string   // Title recomended
	Subtitle string   // Optional
	Cols     string   // Used to set the number of BootStrap colums
	Links    []string // array of links, if we have any
	Button   string   // optional button if we need one
	Text     string   // Just some text

	Content interface{}
}

func NewCard(name string) (c *Card) {
	c = &Card{
		Name: name,
		Cols: "col-4",
	}
	return c
}

func (c *Card) WriteHTML(w http.ResponseWriter, tmpl *template.Template) {
	err := tmpl.ExecuteTemplate(w, "card.html", c)
	if err != nil {
		log.Fatalln(err)
	}
}
