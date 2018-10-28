package moni

import (
	"io"
)

type Object interface {
	io.Reader
	io.Writer
}

type Storage interface {
	UseStore(name string) Storage
	List() (names []string, err error)
	Store(name string, obj interface{}) (err error)
	Fetch(name string, obj interface{}) (err error)
}

type StorageInfo struct {
	Name string
}
