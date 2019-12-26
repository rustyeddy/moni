package main

import (
	"testing"

	"github.com/rustyeddy/store"
)

func TestSiteStorage(t *testing.T) {
	st, err := store.UseFileStore("etc")
	if err != nil {
		t.Fatalf("Failed to open storage %s", err)
	}

	var sitelist []string
	if err = st.ReadObject("sites.json", &sitelist); err != nil {
		t.Fatalf("Failed to read sites.json %v", err)
	}

	if len(sitelist) < 1 {
		t.Errorf("Sites expected > 1 got %d", len(sitelist))
	}
}
