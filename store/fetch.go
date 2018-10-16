package store

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime"
	"net/http"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

// FetchObject unmarshal the contents of the file using the type
// template passed in by otype.  The original Go object will be
// is decoded if desired (e.g. JSON), and a Go object is returned. Nil
// is returned if thier is a problem, such as no object existing.
func (s *Store) FetchObject(name string, otype interface{}) (obj *Object, err error) {
	if obj = s.Get(name); obj == nil {
		return nil, fmt.Errorf("Fetch Object does NOT EXIST %s", name)
	}

	// If our buffer is nil, we will need to fetch the data from the store.
	if obj.Buffer == nil {
		obj.Buffer, err = ioutil.ReadFile(obj.Path)
		if err != nil {
			return nil, fmt.Errorf("  FetchObject failed reading %s -> %v\n", obj.Path, err)
		}
	}
	log.Debugf("  ++ found %d bytes from %s ", len(obj.Buffer), obj.Path)

	// Determine the content type we are dealing with
	ext := filepath.Ext(obj.Path)
	if ext != "" {
		if obj.ContentType = mime.TypeByExtension(ext); obj.ContentType == "" {
			obj.ContentType = http.DetectContentType(obj.Buffer)
		}
	}

	if obj.ContentType == "application/json" {
		if err := json.Unmarshal(obj.Buffer, otype); err != nil {
			return nil, fmt.Errorf("%s: %v", name, err)
		}
	}
	return obj, nil
}
