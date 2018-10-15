/*

Store is a simple library that treats the local filesystem as
an object store.

Storage is an interface that defines a simple Object Storage
interface, allowing alternative implementation to replace or
extend storage from the local filesystem to alternates Object
Stores like s3, Google Cloud Platform, and Digital Ocean.

## Overview

The concept is very simple, the store has containers (directories),
each containter holds objects (files).  Currently we do not support
hiearchies (subdirectories) keeping things simple.

## Containers

Containers can be represented by a directory (container) with
multiple files each representing a single object.

Or a container can be a single file (json, csv, xml, tar.gz.z) that
contains multiple objects.  In this case, the objects are handled
by the application that read or writes the container. Hence with the
file container we do not get involved with the type of "objects" in
the container.

## Objects

Objects can be any type of file "text" or "binary".  If the object
being stored is a Go object, it will be serialized to JSON and
deserialized on read.


## API

http://otto.makerof.site/

- GET /store						# Store meta info
- GET /store/						# list of containers
- GET /store/{container}			# Container meta info
- DEL /store/{container}            # Delete Container
- PUT /store/{container}			# Create container w
- GET /store/{container}/			# List objects in container
- GET /store/{container}/{object}	# Return the object from container
- PUT /store/{container}/{object}   # Add/Replace the container object
- DEL /store/{container}/{object}   # Delete the container object

*/

package store

import (
	"fmt"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

// Storage is a simple CRUDish interface for micro-services that
// may have a local filesystem (VMs) or not (serverless apps),
// or may use cloud bucket storage like S3, GCP or even HTTP.
type Storage interface {
	UseStore(path string) *Store
	ListObjects() (keys []string)
	StoreObject(key string, data interface{}) (err error)
	FetchObject(key string) (obj *Object, err error)
}

// Store is the main structure of this package.  It has a name and
// maintains a path, ObjectIndex and some house keeping private fields
type Store struct {
	Path    string // basedir of this store
	Name    string // the name of the store provider
	Created string
	index
}

// ====================================================================
//                        Store
// ====================================================================

// UseStore creates and returns a new storage container.  If a dir
// already exists, that directory will be used.  If the directory
// does not exist it will be created.
func UseStore(path string) (s *Store, err error) {
	path = filepath.Clean(path)
	s = &Store{
		Path:    path,
		Name:    NameFromPath(path),
		Created: timeStamp(),
		index:   make(index),
	}

	// Determine if we are using an existing directory or need
	// to create a new one.
	// TODO - XXX - Add a permission check.
	if _, err = os.Stat(path); os.IsNotExist(err) {
		// create the path that did not previously exist
		if err := os.MkdirAll(path, 0755); err != nil {
			return nil, fmt.Errorf("mkdir %s: %v", path, err)
		}
	}
	// We already have the store get an index
	s.buildIndex()
	log.Debugln(s.String())
	return s, nil
}

// String a simple summary of our store
func (s *Store) String() string {
	return fmt.Sprintf("name: %s path %s, object count: %d",
		s.Name, s.Path, len(s.Index()))
}

// Index gives us a map indexed by filenames with pointers to
// the corresponding object.  The index needs to be rebuilt as
// the result of any change to the underlying filesystem.  Likewise
// changes to index or any in memory storage will need to be
// flushed to the underlying storage.
func (s *Store) Index() index {
	if s.index == nil {
		s.buildIndex()
	}
	return s.index
}
