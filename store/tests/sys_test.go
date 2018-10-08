package tests

import (
	"os"
	"testing"

	"github.com/rustyeddy/store"
)

var (
	cleanup bool = true
)

type ImportantInfo struct {
	Name   string
	Age    int
	Active bool
}

func getImportantInfo() *ImportantInfo {
	return &ImportantInfo{"Rusty", 102, false}
}

// Create a temporary directory with a random name.  The function
// call will ensure the directory does not already exist
func TestCRUD(t *testing.T) {

	var (
		fs  *store.Store
		obj *store.Object
		err error
	)

	// Get a new empty directory
	testdir := "var/sys-test"
	os.RemoveAll(testdir)
	os.MkdirAll(testdir, 0755)

	// Use the new directory as the store
	if fs, err = store.UseStore(testdir); err != nil {
		t.Error("failed to create store", testdir)
		return
	}

	// Get an object to store
	iinfo := getImportantInfo()

	// Store the object
	if _, err = fs.StoreObject("important1", iinfo); err != nil {
		t.Error(err)
		t.Fail()
	}

	// Check the length of store
	if fs.Len() != 1 {
		t.Errorf("expected Store len (1) got (%d)", fs.Len())
	}

	// Get the store object back, need a pointer to an object
	var info *ImportantInfo

	// Actually Fetch the object.
	if obj, err = fs.FetchObject("important1", &info); err != nil {
		t.Error(err)
		t.Fail()
	}

	// Is it the same object (equivalent)
	if *iinfo != *info {
		t.Errorf("iinfo differs from info")
		t.Logf("\texpected %+v", iinfo)
		t.Logf("\tgot      %+v", info)
	}

	// Remove the object from the store
	key := obj.Name
	if err = fs.DeleteObject(key); err != nil {
		t.Errorf("error removing index %s, %v", key, err)
	}

	// Ensure the object no longer exists
	ex := fs.Exists(key)
	if ex {
		t.Errorf("remove key ~ expected () got key (%s)", key)
	}
}
