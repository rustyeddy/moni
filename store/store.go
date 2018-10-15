package store

/*
	Store is a place to store things
*/

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

// timeStamp returns a timestamp in a modified RFC3339 format,
// basically remove all colons ':' from filename, since they have a
// specific use with Unix pathnames, hence must be escaped when used
// in a filename.
func timeStamp() string {
	ts := time.Now().UTC().Format(time.RFC3339)
	return strings.Replace(ts, ":", "", -1) // get rid of offesnive colons
}

/*
```````````````````````````````````````````````````````````````````````
                     Store and Fetch

    These functions take an empty Go interface, serialize it into a
    the specified format, e.g. JSON, YML, etc. then write the file
    to the underlying container.

```````````````````````````````````````````````````````````````````````
*/

// StoreObject accepts any Go structure, serialize it as a JSON object.
// The JSON object will then be written to disk encapsulated as an "Object".
// The Object contains some meta data about the original object, including
// it's type to help the Application be able to deserialize and use the
// Object with out implicit knowledge of the objects structure.
func (s *Store) StoreObject(name string, data interface{}) (obj *Object, err error) {

	// Do NOT allow '/' characters in string
	if strings.Index(name, "/") > -1 {
		return nil, fmt.Errorf("illegal char '/' used for index %s", name)
	}

	stobj := data                    // turn an interface into an object
	jbuf, err := json.Marshal(stobj) // JSONify the "object" param
	if err != nil {
		return nil, fmt.Errorf("JSON marshal failed %s -> %v", name, err)
	}

	// log.Debug("  storing data :", string(jbuf[0:40]))
	obj = ObjectFromBytes(jbuf) // obj will not be nil
	if obj == nil {
		log.Fatalln("SUCKATARI")
		return
	}
	obj.Store = s // back pointer to store
	obj.Name = name
	obj.Path = s.Path + "/" + name + ".json"

	// Now write to the file
	err = ioutil.WriteFile(obj.Path, jbuf, 0644)
	if err != nil {
		return nil, fmt.Errorf("  Store.Object write failed %s -> %v", obj.Path, err)
	}
	s.Set(name, obj)
	return obj, nil
}

// DeleteObject does just that, it removes the object from the store.
// Meaning it removes the object from the disk
func (s *Store) DeleteObject(name string) error {
	var (
		obj *Object
	)
	if obj = s.Get(name); obj == nil {
		return fmt.Errorf("%s NOT FOUND", name)
	}

	// The object must be removed from the filesystem first ...
	if obj.Path == "" {
		return fmt.Errorf("path is nil, should never happen %s", name)
	}
	if err := os.Remove(obj.Path); err != nil {
		return fmt.Errorf("Remove path %s error %v", obj.Path, err)
	}

	// Now remove form the index.
	delete(s.index, name)
	return nil
}

// =======================================================================
// Index returns a map of item names and full paths
// =======================================================================

// Index will scan the store directory for objects (files) creating a
// map of pointers to the Objects indexed by the object name (file
// name less the path and extension)
func (s *Store) buildIndex() index {
	// Now build the index if we don't have one
	s.Path = filepath.Clean(s.Path) // Cleanse our path
	pattern := s.Path + "/*"
	s.indexPaths(pattern)
	return s.index
}

// indexPaths will create a map of *File created from fullpaths indexed by
// the filename (less extension).
func (s *Store) indexPaths(pattern string) (err error) {
	var (
		paths []string
	)

	if paths, err = filepath.Glob(pattern); err != nil || paths == nil {
		return fmt.Errorf("no files to index %s %v", pattern, err)
	}

	// Create room in the index for the paths
	if s.index == nil {
		s.index = make(index, len(paths))
	}

	// for the range of paths
	for _, p := range paths {
		var (
			fi  os.FileInfo
			err error
		)

		// We only want to index regular files, Lstat will help use determine
		if fi, err = os.Lstat(p); err != nil {
			log.Warningln("Lstat error ", p, err) // TODO: Append to a buffer
			continue
		}

		// We want to log this incase a time comes later if we do
		// care about non-regular files
		if !fi.Mode().IsRegular() {

			// Should we complain about directories?
			log.Debugf("  ignore (non regular file) %s", p)
			continue
		}

		var obj *Object
		if obj, err = ObjectFromPath(p); err != nil {
			log.Errorln(p, err)
		}
		obj.Store = s

		// attach the object to the index
		s.Set(obj.Name, obj)
	}
	return nil
}

// Count returns the number of items in Store
func (s *Store) Count() int {
	return len(s.Index())
}

// ==============================================================================

func (s *Store) Shutdown() {
	// Flush in memory data
	// Shutdown connections
}
