package moni

import (
	"net/http"

	"github.com/rustyeddy/moni/app"
)

// AppHandler will compose a response to the request
// and in the process most likely need to gather a few peices
// of information, put them together in the right order and
// send back to the caller
func AppHandler(w http.ResponseWriter, r *http.Request) {

	pb := app.NewPageBuilder()
	pb.Title = "Application Interface"
	pb.PageData.Name = "Rusty Eddy"
	pb.Frag = r.URL.Fragment

	card := &app.Card{
		Title: "Card is a Good ONe",
		Text:  "Hey man, I got some shit to say",
		Links: []*app.Link{
			&app.Link{"Clowd Ops", "http://clowdops.net"},
			&app.Link{"USC", "http://usc.edu"},
		},
	}
	pb.AddCard(w, card)

	ncard := *card
	card.Title = "A new version of the Card!"
	card.Text = "If you think I'm sexy come on sugar let me know"
	pb.AddCard(w, &ncard)

	pb.Assemble(w, pb.TemplateName)
}
