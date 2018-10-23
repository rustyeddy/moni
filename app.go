package moni

import (
	"net/http"
	"strings"

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

	c1 := GetCrawlCard()
	pb.AddCard(w, c1)

	c2 := GetCard2()
	pb.AddCard(w, c2)

	c3 := GetCard3()
	pb.AddCard(w, c3)
	pb.Assemble(w, "index.html")
}

func GetCrawlCard() *app.Card {
	card := app.Card{
		Title: "Recent Crawls ...",
		Cols:  "col-6",
		Links: []*app.Link{
			&app.Link{"Clowd Ops", "http://clowdops.net"},
			&app.Link{"USC", "http://usc.edu"},
		},
	}
	card.Text = strings.Join(GetCrawls(), " <br/>\n\t")
	return &card
}

func GetCard2() *app.Card {
	card := app.Card{
		Title: "This is a NEW Card",
		Text:  "Hey new!",
		Cols:  "col-6",
		Links: []*app.Link{
			&app.Link{"Clowd Ops", "http://clowdops.net"},
			&app.Link{"USC", "http://usc.edu"},
		},
	}
	return &card
}

func GetCard3() *app.Card {
	card := app.Card{
		Title: "This is a NEW Card",
		Text:  "Hey man, I am brand new!",
		Cols:  "col-3",
		Links: []*app.Link{
			&app.Link{"Clowd Ops", "http://clowdops.net"},
			&app.Link{"USC", "http://usc.edu"},
		},
	}
	return &card
}
