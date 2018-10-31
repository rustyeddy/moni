package moni

import (
	"fmt"
	"net/http"
)

func HostSiteHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, map[string]string{
		"host": "localhost",
		"ip":   "127.0.0.1",
	})
}
