package store

import (
	"fmt"
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// Object represents a single object contained by Store.  It has a
// name unique to this store.  The object has a File representing the
// objects location on disk and Content representing the actual data
// that belongs to this object, in the format specified by the content
type Object struct {
	Name        string // the name as registered by Store
	Path        string // full path in filesystem
	ContentType string // mime/type
	Type        string

	*Store        // point to our storage
	Buffer []byte // raw data

	os.FileInfo
}

// ObjectFromPath will create a new object and prefill
// with path, name and mimetype derived from path extension.
func ObjectFromPath(path string) (obj *Object, err error) {
	var fi os.FileInfo
	if fi, err = os.Stat(path); os.IsPermission(err) {
		fmt.Println("2")
		return nil, fmt.Errorf("path %s: %v", path, err)
	} else if os.IsNotExist(err) {
		return nil, fmt.Errorf("%s: %v", path, err)
	} else if err != nil {
		return nil, err
	}

	ext := filepath.Ext(path)
	obj = &Object{
		Name:        NameFromPath(path),
		Path:        filepath.Clean(path),
		ContentType: mime.TypeByExtension(ext),
		Buffer:      nil,
		FileInfo:    fi,
	}
	if obj.ContentType == "" {
		obj.ContentType = "application/octet-stream"
	}
	return obj, nil
}

// ContentFromBytes will create a new Content initialized with appropriate
// data.  We will use http.DetectContentType to determine the type of
// content.
func ObjectFromBytes(buf []byte) (obj *Object) {
	obj = &Object{
		Buffer:      buf,
		ContentType: http.DetectContentType(buf),
	}
	return obj
}

// Size returns the size of the bytes in the buffer
func (o *Object) Size() int64 {
	if o.FileInfo != nil {
		return o.FileInfo.Size()
	}
	return -1
}

// the name of the object is the filename excluding the extension
// and path.
func NameFromPath(path string) string {
	_, fname := filepath.Split(path)
	flen := len(fname) - len(filepath.Ext(fname))
	return fname[0:flen] // return less the .ext
}

// Compare two Objects, basically, if they point at the same file the
// will be considered equal regardless of the runtime state (if buffer
// is cached or verbosity)
func (o *Object) Compare(obj *Object) bool {
	if strings.Compare(obj.Name, o.Name) != 0 {
		return false
	}
	if strings.Compare(obj.Path, o.Path) != 0 {
		return false
	}
	if strings.Compare(obj.ContentType, o.ContentType) != 0 {
		return false
	}
	return true
}

// Bytes returns a pointer to the bytes
func (o *Object) Bytes() []byte {
	return o.Buffer
}

// Write data to the object
func (o *Object) Write(b []byte) (n int, err error) {
	err = ioutil.WriteFile(o.Path, b, 0644)
	return len(b), err
}

// Read data from the file (append or rewrite)
func (o *Object) Read(b []byte) (n int, err error) {
	b, err = ioutil.ReadFile(o.Path)
	return len(b), err
}
