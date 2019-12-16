package main

import (
	"net/url"

	log "github.com/sirupsen/logrus"
)

type Site struct {
	*url.URL
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
