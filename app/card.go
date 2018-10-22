package ui

import "log"

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
			c.Links = v
		case "Image":
			c.Image.Src = v
		case "Alt":
			c.Image.Alt = v
		default:
			log.Errorf("unexpected card part %s, ignoring .. %s with value ", k)
		}
	}
	return c
}

// HTML ~ Hmmm do we want to generate html directly?  Hmm. probably
func (c *Card) HTML() (h string) {

}
