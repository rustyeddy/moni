package moni

import (
	"github.com/rustyeddy/store"
)

var (
	storage *store.Store
)

func GetStorage() (st *store.Store) {
	var err error

	dir := "/srv/moni"
	if config != nil && config.ConfigFile != "" {
		dir = config.ConfigFile
	}
	if st, err = store.UseStore(dir); err != nil {
		return st
	}
	log.Fatalf("failed to get storage %s err %v ", dir, err)
	return st
}
