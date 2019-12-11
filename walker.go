package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	log "github.com/sirupsen/logrus"
)

// Walk response manage the walk of this particular host
type Walker struct {
	URLstr string `json:"url"`
	*Page  `json:"page"`
	io.Writer
}

func (w *Walker) Write(b []byte) (n int, err error) {
	fmt.Fprintf(os.Stdout, "%v\n", b)
	return len(b), nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	u := r.URL.Query().Get("url")
	processURL(u, w)
	log.Println("handler: sending urlstr to urlChan", u)
}

// Scrubber spins listening to the urlChan acceptring URLs that
// need to be scrubbed.
func Scrubber(c chan Walker, d chan bool) {
	var page *Page
	log.Infoln("Starting the Scrubber ..")

	for {
		log.Infoln("\turlChan waiting for an URL")
		select {
		case walker := <-c:
			log.Infof("\tgot urlstring %s\n", walker.URLstr)
			if u := scrubURL(walker.URLstr); u != nil {
				if page = GetPage(*u); page != nil {
					log.Infof("got page: %+v - let's walk...\n", page)
					go page.Walk(walker.Writer)
				}
			}

			// Add this if we want to timeout from an interactive run
			// case <-time.After(2 * time.Second):
			// 	d <- true
		}
	}
}

func scrubURL(urlstr string) (u *url.URL) {
	var err error

	log.Infoln("scrubURL with ", urlstr)

	u, err = url.Parse(urlstr)
	errPanic(err)

	if u.Scheme == "" {
		u.Scheme = "http"
	}

	if u, err = url.Parse(u.String()); err != nil {
		return nil
	}

	// if this hostname exists in the acl set as false,
	// we will just return
	if f, ex := acl[u.Host]; ex && f == false {
		return nil
	}
	return u
}
