package ui

// CardOptions
type CardOptions map[string]string

// Card represents a card from our Bootstrap Dashboard
// CSS Template
type Card struct {
	*Image
	Title string  // header
	Text  string  // body
	Links []*Link //

	*html.Template
}

type Link struct {
	Text string
	Link string
}

// NewCard creates a new card
func NewCard(config map[string]string) {

}
