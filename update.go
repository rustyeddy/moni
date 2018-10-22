package moni

import (
	"fmt"
	"net/http"
)

func UpdateHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("%+v\n", r)
	fmt.Fprintln(w, "got it")
}
