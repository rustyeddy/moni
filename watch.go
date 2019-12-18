package main

import (
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

// watchSites will run in an endless loop separated by a sleep
func watchSites(wg *sync.WaitGroup) {
	defer wg.Done()

	sitelist = GetSites()
	log.Infof("Starting watchSites: %v", sitelist)
	for {
		for _, s := range sitelist {

			// TODO: This function will print to std output, it we
			// should be passing an io.Writer rather than nil.  We are
			// currently switching between stdout and the
			// http.ResponseWriter
			processURL(s, nil)
		}

		stime := time.Duration(config.Wait) * time.Minute
		log.Infof("sleeping %s", stime)
		time.Sleep(stime)
	}
}
