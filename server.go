package inv

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type Server struct {
	Addrport string // config option
	Pubdir   string // config option

	*mux.Router
	*log.Logger
}

var (
	r *mux.Router
)

/*
   Use the runtime/trace and net/http/pprof packages
*/

// HTTPServer will register routes, open connection and listen for incoming
func StartServer(addrport string, done chan<- bool) {
	if addrport == "" {
		log.Fatalf("Server must have a port to listen with")
	}
	r = mux.NewRouter()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Home Handler")
	})
	r.HandleFunc("/crawl/{url}", HandleCrawl)
	// r.Handle("/", http.FileServer(http.Dir("./pub")))

	log.Infoln("listening on ", addrport)
	err := http.ListenAndServe(addrport, r)
	if err != nil {
		log.Errorf("Server terminated error %v", err)
	}
	done <- true
}

// Send back diagnostic info
func handleInfo(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Handling an HTML request ")

	info := make(map[string]string)
	info["Routes"] = "/, /info, /walk/{url}"
	msg, err := json.Marshal(&info)
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	fmt.Fprint(w, msg)
}
