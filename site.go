package main

import (
	"net/url"

	log "github.com/sirupsen/logrus"
)

type Site struct {
	*url.URL
	Pages map[string]*Page
}

func GetSites() []string {
	if sitelist == nil || len(sitelist) < 1 {
		sitelist = make([]string, 1)
		if err := storage.ReadObject("sites.json", &sitelist); err != nil {
			log.Errorf("Storage failed to read sites.json: %v", err)
		}
	}
	return sitelist
}

// SaveSites saves the sites structure.
func SaveSites() (err error) {
	if sitelist != nil && len(sitelist) > 0 {
		if err = storage.Save("sites.json", &sitelist); err != nil {
			log.Errorf("Storage Save failed for sites.json %v", err)
		}
	}
	return err
}

func AddSite(url string) []string {

	// first make sure we are not added the same website twice
	// XXX: this will.
	sitelist = append(sitelist, url)
	return sitelist
}
