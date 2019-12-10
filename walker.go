package main

import (
	"net/http"
	"net/url"

	log "github.com/sirupsen/logrus"
)

func handler(w http.ResponseWriter, r *http.Request) {
	urlstr := r.URL.Query().Get("url")
	if urlstr == "" {
		log.Println("missing URL argument")
		return
	}

	log.Println("handler: sending urlstr to urlChan", urlstr)
	urlChan <- urlstr
}

// Scrubber spins listening to the urlChan acceptring URLs that
// need to be scrubbed.
func Scrubber(c chan string, d chan bool) {
	var page *Page
	log.Infoln("Starting the Scrubber ..")

	for {
		log.Infoln("\turlChan waiting for an URL")
		select {
		case urlstr := <-c:
			log.Infof("\tgot urlstring %s\n", urlstr)
			if u := scrubURL(urlstr); u != nil {
				if page = GetPage(*u); page != nil {
					log.Infof("got page: %+v - let's walk...\n", page)
					page.Walk()
				}
			}
		
		case <-time.After(1 * time.Second):
					
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
