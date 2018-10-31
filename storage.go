package moni

import (
	"encoding/json"
	"io"
	"io/ioutil"
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

func UseStore(dir string) (s *Store) {
	return &Store{
		Basedir:     dir,
		ContentType: "application/json",
		Ext:         ".json",
	}
}

func (s *Store) Put(name string, obj interface{}) (err error) {
	var buf []byte

	switch s.ContentType {
	case "application/json":
		if buf, err = json.Marshal(obj); err != nil {
			IfErrorFatal(err, "marshaling json "+name)
		}
	default:
		panic("did not expect this")
	}

	// Write the file to disk
	if err = ioutil.WriteFile(name, buf, 0755); err != nil {
		IfErrorFatal(err, "writing buffer "+name)
	}
	return err
}

func (s *Store) Get(name string, obj interface{}) (err error) {
	var buf []byte

	if buf, err = ioutil.ReadFile(name); err != nil {
		IfErrorFatal(err, "s.Get readfile"+name)
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
