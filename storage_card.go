package moni

import "github.com/rustyeddy/store"

// SitesCard provides a card with list of the sites we are managing
type StorageCard struct {
	*Card
	*store.Store
}

func NewStorageCard() *StorageCard {
	return &StorageCard{
		Card:  NewCard("Storage"),
		Store: GetStorage(),
	}
}
