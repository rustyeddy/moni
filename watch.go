package main

import (
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

// watchSites will run in an endless loop separated by a sleep
func watchSites(wg *sync.WaitGroup) {
	defer wg.Done()

	sites = GetSites()
	log.Infof("Starting watchSites: %v", sites)
	for {
		for _, s := range sites {

			// TODO: This function will print to std output, it we
			// should be passing an io.Writer rather than nil.  We are
			// currently switching between stdout and the
			// http.ResponseWriter
			processURL(s, nil)
		}
		log.Infoln()
		time.Sleep(5 * time.Minute)
	}
}
