package moni

import (
	"github.com/rustyeddy/store"
	log "github.com/sirupsen/logrus"
)

var (
	moniStore *store.Store
)

func GetStorage() *store.Store {
	var err error

	if moniStore == nil {
		if moniStore, err = store.UseStore(config.StoreDir); err != nil {
			log.Fatalf("Unable to access our store")
		}
	}
	return moniStore
}
