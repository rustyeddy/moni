package moni

import (
	"net/http"
)

// AppHandler will compose a response to the request
// and in the process most likely need to gather a few peices
// of information, put them together in the right order and
// send back to the caller
func AppHandler(w http.ResponseWriter, r *http.Request) {

	pb := NewPageBuilder()
	pb.Title = "Application Interface"
	pb.Name = "Rusty Eddy"
	pb.Frag = r.URL.Fragment

	pb.Assemble(w, "index.html")
}
