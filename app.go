package main

import (
	"fmt"
	"net/http"
)

func AppHandler(w http.ResponseWriter, r *http.Request) {
	// Serve up the website slash app
	fmt.Println("Hello, World")
}
