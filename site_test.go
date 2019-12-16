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

	if err = st.ReadObject("sites.json", &sites); err != nil {
		t.Fatalf("Failed to read sites.json %v", err)
	}

	if len(sites) < 1 {
		t.Errorf("Sites expected > 1 got %d", len(sites))
	}
}
