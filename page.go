package main

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

// ===================================================================
type Page struct {
	URL string

	Content []byte
	Links   map[string]*Page
	Ignored map[string]int

	StatusCode int

	LastCrawled time.Time
	Start       time.Time
	Finish      time.Time
	crawl       bool
	Err         error
}

var (
	pages Pagemap
)

// Pagemap
// ********************************************************************
type Pagemap map[string]*Page

// String will represent the Page
// ====================================================================
func (p *Page) String() string {
	str := fmt.Sprintf("%s: lastcrawled: %s,  duration: %v links: %d ignored: %d\n", p.URL, p.LastCrawled, p.Finish, len(p.Links), len(p.Ignored))
	return str
}

func GetPages() Pagemap {
	if pages == nil {
		pages = make(Pagemap)

		st := getStorage()
		if _, err := st.FetchObject("pages", &pages); err != nil {
			log.Debugf("Empty pages %v, creating ...", err)
			// TODO ~ make sure the error is NOT found

			pages = make(Pagemap)
			_, err := st.StoreObject("pages", pages)
			if err != nil {
				log.Errorf("Store: failed to create pages: %v ", err)
				return pages
			}
		}
	}
	return pages
}

func savePagemap() error {
	st := getStorage()
	if _, err := st.StoreObject("pages", pages); err != nil {
		log.Errorf("failed to save page map %v", err)
		return err
	}
	return nil
}

// GetPage will sanitize the url, either find or create the
// corresponding page structure.  If the URL is deep, we also
// find the corresponding site structure.
func PageFromURL(ustr string) (pi *Page) {
	var ex bool
	if pi, ex = pages[ustr]; !ex {
		pi = &Page{
			URL:     ustr,
			Links:   make(map[string]*Page),
			Ignored: make(map[string]int),
			crawl:   true,
		}
		pages[ustr] = pi
	}
	return pi
}

func (pm Pagemap) Get(url string) (p *Page) {
	if p, e := pages[url]; e {
		return p
	}
	return nil
}

func (pm Pagemap) Exists(url string) bool {
	if p := pm.Get(url); p != nil {
		return true
	}
	return false
}

func (pm Pagemap) Set(url string, p *Page) {
	pages[url] = p
}

func PageListHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, pages)
}

func PageIdHandler(w http.ResponseWriter, r *http.Request) {

	// Get a couple vars ready for later
	var out interface{}
	var err error

	// Get the url from the request and extract the storage
	// index (name) of the corresponding object.
	url := urlFromRequest(r)
	name := NameFromURL(url)
	out = "No Output"

	if page := pages.Get(url); page != nil {
		switch r.Method {
		case "GET":
			out = page
		case "PUT", "POST":
			log.Infoln("overwriting ", name)
			out = "done"
		case "DELETE":
			delete(pages, name)
			out = "done"
		}
	} else {
		switch r.Method {
		case "GET", "DELETE":
			// Nothing to get or delete
			err = errors.New("object not found " + name)
		case "PUT", "POST":
			st := getStorage()
			if _, err := st.StoreObject(name, page); err != nil {
				err = fmt.Errorf("page %s error %v", url, err)
			} else {
				out = `{"msg": "done"}`
			}
		}
	}
	if err != nil {
		JSONError(w, err)
	}
	writeJSON(w, out)
}
