package moni

import (
	"log"

	"github.com/rustyeddy/store"
)

var (
	storage *store.Store
)

func init() {
	var err error
	if storage, err = store.UseStore(config.Storedir); err != nil {
		log.Fatalf("failed to get storage %s ", config.Storedir)
	}
}
