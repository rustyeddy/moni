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

	card := app.Card{
		Title: "Card is a Good ONe",
		Text:  "Hey man, I got some shit to say",
		Cols:  "col-6",
		Links: []*app.Link{
			&app.Link{"Clowd Ops", "http://clowdops.net"},
			&app.Link{"USC", "http://usc.edu"},
		},
	}
	pb.AddCard(w, &card)

	ncard := card
	ncard.Title = "A new version of the Card!"
	ncard.Text = "If you think I'm sexy come on sugar let me know"
	ncard.Cols = "col-3"
	pb.AddCard(w, &ncard)

	zcard := ncard
	zcard.Title = "I AM Right!"
	zcard.Text = "Just ask me."
	zcard.Cols = "col-3"
	pb.AddCard(w, &zcard)

	pb.Assemble(w, pb.TemplateName)
}
