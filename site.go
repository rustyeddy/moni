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
	if sites == nil || len(sites) < 1 {
		sites = make([]string, 1)
		if err := storage.ReadObject("sites.json", &sites); err != nil {
			log.Errorf("Storage failed to read sites.json: %v", err)
		}
	}
	return sites
}

// SaveSites saves the sites structure.
func SaveSites() (err error) {
	if sites != nil && len(sites) > 0 {
		if err = storage.Save("sites.json", &sites); err != nil {
			log.Errorf("Storage Save failed for sites.json %v", err)
		}
	}
	return err
}

func AddSite(url string) []string {
	sites = append(sites, url)
	return sites
}
