package main

import "github.com/rustyeddy/store"

func setupStorage() {
	if storage, err = store.UseFileStore("."); err != nil || storage == nil {
		errFatal(err, "failed to useFileStore ")
	}
}
