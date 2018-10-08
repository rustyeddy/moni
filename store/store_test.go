package store

import (
	"os"
	"testing"
)

var (
	st       *Store
	testpath string = "../var/store-tests"
)

type tKV struct {
	K string
	V string
}

type tKeyVal struct {
	Key string
	Val interface{}
}

func ResetTestDir() string {
	os.RemoveAll(testpath)
	os.MkdirAll(testpath, 0755)
	return testpath
}

// BadStoreTest will ensure that the library reports an error when we
// ask it to use an illegal path for storage.
func TestBadStore(t *testing.T) {
	badpath := "/badpath/dont/be/root/"
	if st, err := UseStore(badpath); err == nil {
		t.Errorf("path should not have been created")
	} else if st != nil {
		t.Errorf("badpath expected (err) got (%v) ", st)
	}
}

// NewStoreTest will test using the newstore
func TestNewStore(t *testing.T) {
	var (
		err error
		st  *Store
	)
	v1 := []tKV{
		{"one", "1"},
		{"two", "2"},
		{"three", "3"},
	}
	v2 := []tKeyVal{
		{"1", 1},
		{"2", 2},
		{"3", 3},
	}

	testpath = ResetTestDir()
	if st, err = UseStore(testpath); err != nil {
		t.Errorf("expected (store) got error (%v)", err)
	}

	var oa, ob *Object
	if _, err = st.StoreObject("idx1", v1); err != nil {
		t.Error("idx1", err)
	}

	if _, err = st.StoreObject("index2", v2); err != nil {
		t.Error("index2", err)
	}

	// lets check that we have two items in the index
	if st.Len() != 2 {
		t.Errorf("index expected (2) got (%d) ", len(st.Index()))
	}

	for n, o := range st.Index() {
		switch n {
		case "idx1":
			oa = o
		case "index2":
			ob = o
		default:
			t.Errorf("expected index got (%s) ", o.Name)
		}
	}

	if oa == nil {
		t.Errorf("expected (idx1) got (%s)", "")
	}
	if ob == nil {
		t.Errorf("expected (idx2) got (%s)", "")
	}
}

func TestExistingStore(t *testing.T) {
	var (
		st     *Store
		kv     []tKV
		keyval []tKeyVal
		err    error
	)

	if st, err = UseStore(testpath); err != nil {
		t.Error(testpath, err)
	}

	if st.Len() != 2 {
		t.Errorf("index len expected (2) got (%d) ", st.Len())
	}

	if _, err = st.FetchObject("idx1", &kv); err != nil {
		t.Errorf("expected idx1 got (%v)", err)
	}

	if len(kv) != 3 {
		t.Errorf("expected kv len (3) got (%d) %+v", len(kv), kv)
	}

	for _, v := range kv {
		failed := true
		switch v.K {
		case "1":
			failed = (v.V == "one")
		case "2":
			failed = (v.V == "two")
		case "3":
			failed = (v.V == "three")
		default:
			failed = false
		}
		if failed {
			t.Error("failed", v.K, v.V)
		}
	}

	if _, err = st.FetchObject("index2", &keyval); err != nil {
		t.Errorf("expected index2 got (%v)", err)
	}

	if len(keyval) != 3 {
		t.Errorf("expected keyval len (%d) got (%d)", 2, len(keyval))
	}
}

func TestDeleteObject(t *testing.T) {
	var (
		st  *Store
		err error
	)

	if st, err = UseStore(testpath); err != nil {
		t.Error(err)
	}

	if st.Count() != 2 {
		t.Errorf("expected (2) objects got (%d) ", st.Count())
	}

	idx := "index2"
	if err = st.DeleteObject(idx); err != nil {
		t.Error("failed to delete ", idx)
	}

	if st.Count() != 1 {
		t.Errorf("count expected (1) got (%d)", st.Count())
	}
}

func TestMain(m *testing.M) {

	// Setup test directory
	testdir := ResetTestDir()

	// Clean up
	defer os.RemoveAll(testdir)

	// Get results from test run
	results := m.Run()

	// Tear down suff...

	// Exist with the results from the test run
	os.Exit(results)
}
