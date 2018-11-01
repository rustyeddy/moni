package moni

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"path/filepath"
)

type Object interface {
	io.Reader
	io.Writer
}

type Storage interface {
	UseStore(name string) Storage
	Index() (names []string, err error)
	Store(name string, obj interface{}) (err error)
	Fetch(name string, obj interface{}) (err error)
	Close() (err error)
}

type Store struct {
	Basedir     string
	ContentType string
	Ext         string
}

var (
	store *Store
)

func GetStore() *Store {
	if store == nil {
		store = UseStore(config.Storedir)
	}
	return store
}

func UseStore(dir string) (s *Store) {
	return &Store{
		Basedir:     dir,
		ContentType: "application/json",
		Ext:         ".json",
	}
}

func (s *Store) Glob(pattern string) []string {
	matches, err := filepath.Glob(pattern)
	IfErrorFatal(err, "Glob")
	return matches
}

func (s *Store) pathFromName(name string) string {
	return s.Basedir + "/" + name
}

func (s *Store) Put(name string, obj interface{}) (err error) {
	var buf []byte

	switch s.ContentType {
	case "application/json":
		if buf, err = json.Marshal(obj); err != nil {
			IfErrorFatal(err, "marshaling json "+name)
		}
	default:
		log.Fatalf("did not expect ContentType %s", s.ContentType)
	}

	// Write the file to disk
	path := s.pathFromName(name)
	if err = ioutil.WriteFile(path, buf, 0755); err != nil {
		IfErrorFatal(err, "writing buffer "+path)
	}
	return err
}

func (s *Store) Get(name string, obj interface{}) (err error) {
	var buf []byte

	path := s.pathFromName(name)
	if buf, err = ioutil.ReadFile(path); err != nil {
		return fmt.Errorf("read index %s failed %v", path, err)
	}

	switch s.ContentType {
	case "application/json":
		if err = json.Unmarshal(buf, obj); err != nil {
			IfErrorFatal(err, "marshaling json "+name)
		}
	default:
		panic("did not expect this")
	}
	return err
}
