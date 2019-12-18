package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

// spaHandler implements the http.Handler interface, so we can use it
// to respond to HTTP requests. The path to the static directory and
// path to the index file within that static directory are used to
// serve the SPA in the given static directory.
type spaHandler struct {
	staticPath string
	indexPath  string
}

// ServeHTTP inspects the URL path to locate a file within the static dir
// on the SPA handler. If a file is found, it will be served. If not, the
// file located at the index path on the SPA handler will be served. This
// is suitable behavior for serving an SPA (single page application).
func (h spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	log.Infof("HTML Server called %+v", r.URL)

	// get the absolute path to prevent directory traversal
	path, err := filepath.Abs(r.URL.Path)
	if err != nil {
		// if we failed to get the absolute path respond with a 400 bad request
		// and stop
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// prepend the path with the path to the static directory
	path = filepath.Join(h.staticPath, path)
	if path == "pub/favicon.ico" || path == "/favicon.ico" {
		return
	}

	// check whether a file exists at the given path
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		// file does not exist, serve index.html
		http.ServeFile(w, r, filepath.Join(h.staticPath, h.indexPath))
		return
	} else if err != nil {
		// if we got an error (that wasn't that the file doesn't exist) stating the
		// file, return a 500 internal server error and stop
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// otherwise, use http.FileServer to serve the static dir
	http.FileServer(http.Dir(h.staticPath)).ServeHTTP(w, r)
}

func doRouter(dir string, wg *sync.WaitGroup) {
	router := mux.NewRouter()
	defer wg.Done()

	router.HandleFunc("/api/config", getConfig)
	router.HandleFunc("/api/config/{key}/{val}", setConfig)
	router.HandleFunc("/api/health", getHealth)
	router.HandleFunc("/api/sites", getSites)
	router.HandleFunc("/api/site/{url}", getSite)
	router.HandleFunc("/api/pages", getPages)

	spa := spaHandler{staticPath: dir, indexPath: "index.html"}
	router.PathPrefix("/").Handler(spa)

	srv := &http.Server{
		Handler: router,
		Addr:    config.Addrport,

		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	err = srv.ListenAndServe()
}

func getHealth(w http.ResponseWriter, r *http.Request) {
	// an example API handler
	json.NewEncoder(w).Encode(map[string]bool{"ok": true})
}

func getConfig(w http.ResponseWriter, r *http.Request) {
	// an example API handler
	json.NewEncoder(w).Encode(config)
}

func setConfig(w http.ResponseWriter, r *http.Request) {
	var err error
	vars := mux.Vars(r)
	key := vars["key"]
	val := vars["val"]
	if r.Method != "POST" && r.Method != "PUT" {
		fmt.Fprintln(w, "Request has bad method %s Bad Form", r.Method)
		return
	}

	switch key {
	case "wait":
		if config.Wait, err = strconv.Atoi(val); err != nil {
			log.Errorf("failed to set configuration %v", err)
			fmt.Fprintln(w, "Error Bad Form ~> ParseForm()")
			return
		}
	}
	json.NewEncoder(w).Encode(config)
}

func getSites(w http.ResponseWriter, r *http.Request) {
	sr := struct {
		Sites []string
	}{
		GetSites(),
	}
	json.NewEncoder(w).Encode(sr)
}

func getSite(w http.ResponseWriter, r *http.Request) {
	var url string

	vars := mux.Vars(r)
	if url = vars["url"]; url == "" {
		fmt.Fprintln(w, "Bad Form ~> ParseForm()")
		return
	}

	log.Infof("Adding %s to sitelist and walking", url)
	AddSite(url)

	// Process the site since it is new, it will return with
	// json.NewEncoder(w.w).Encode(pg)
	processURL(url, w)
}

func getPages(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("PAGES: %+v\n", pages)
	json.NewEncoder(w).Encode(&pages)
}
