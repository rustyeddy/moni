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

// Storage is the interface that any
type Storage interface {
	UseStore(path string) *Store
	ListObjects() (keys []string)
	StoreObject(key string, data interface{}) (err error)
	FetchObject(key string) (obj *Object, err error)
}
