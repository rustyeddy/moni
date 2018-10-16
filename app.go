package main

import (
	"fmt"
	"net/http"
)

func AppHandler(w http.ResponseWriter, r *http.Request) {
	// Serve up the website slash app
	w.Header().Add("Content-Type", "text/html")
	fmt.Fprintln(w, "<h1>Hello, World</h1>")
}
