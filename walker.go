package main

import (
	"fmt"
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

	log.Println("handler: ", urlstr)
	scrubURLs([]string{urlstr})
}

func scrubURLs(urls []string) {
	// walk the command line arguments treating them as URLs
	for _, baseURL := range urls {
		var page *Page

		if url := scrubURL(baseURL); url != nil {
			if page = GetPage(url.String()); page != nil {
				page.Walk()
			}
		}
	}
	storage.Save("config2.json", config)
	storage.Save("pages.json", pages)
	fmt.Println("The end...")
}

func scrubURL(urlstr string) (u *url.URL) {
	var err error

	u, err = url.Parse(urlstr)
	errPanic(err)
	if u.Scheme == "" {
		u.Scheme = "http"
	}

	u, err = url.Parse(u.String())

	// if this hostname exists in the acl set as false,
	// we will just return
	if f, ex := acl[u.Hostname()]; ex && f == false {
		return nil
	}

	// This is a little risky
	acl[u.Hostname()] = true
	return u
}
