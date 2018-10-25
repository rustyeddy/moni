package moni

import "github.com/rustyeddy/store"

// SitesCard provides a card with list of the sites we are managing
type StorageCard struct {
	*Card
	*store.Store
}

// Create some useful information from storage
func NewStorageCard() *StorageCard {
	sc := &StorageCard{
		Card:  NewCard("Storage"),
		Store: GetStorage(),
	}

	return sc
}
