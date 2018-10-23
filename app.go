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

	pb := NewPageBuilder()
	card := &app.Card{
		Image: nil,
		Title: "Card is a Good ONe",
		Text:  "Hey man, I got some shit to say",
		Links: []*Link{
			&Link{"Clowd Ops", "http://clowdops.net"},
			&Link{"USC", "http://usc.edu"},
		},
	}

	pb.AddCard(card)

	d := app.AppData{
		Title: "Application Interface",
		Name:  "Rusty Eddy",
		Frag:  r.URL.Fragment,
		Card:  card,
	}
	pb.Assemble(w, pb.TemplateName, &d)
}
