package main

import (
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

// watchSites will run in an endless loop separated by a sleep
func watchSites(sites []string, wg *sync.WaitGroup) {
	defer wg.Done()

	log.Infof("Starting watchSites: %v", sites)
	for {
		for _, s := range sites {
			processURL(s, nil)
		}
		log.Infoln()
		time.Sleep(5 * time.Minute)
	}

}
