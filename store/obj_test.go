package store

import (
	"testing"
)

var (
	objPaths []string
)

func init() {
	objPaths = []string{".store/idx1.json", ".store/index2.json", "/bad/badpath.json"}
}

func isNil(obj interface{}) bool {
	return obj == nil
}

// Make Sure Objects cabe be created from paths
func TestObjectFromPath(t *testing.T) {
	var err error
	expectedNilObj := []bool{false, false, true}
	expectedNames := []string{"idx1", "index2", ""}
	expectedPaths := []string{".store/idx1.json", ".store/index2.json", ""}
	expectedContentTypes := []string{"application/json", "application/json", ""}

	var obj *Object
	for i, path := range objPaths {

		if obj, err = ObjectFromPath(path); err != nil {
			continue
		}

		isnil := (obj == nil)
		if expectedNilObj[i] != isnil {
			t.Errorf("obj %d expected (%t) got (%t)", i, expectedNilObj[i], isnil)
		}

		name := NameFromPath(path)
		if name != expectedNames[i] {
			t.Errorf("obj %d names expected (%s) got (%s)", i, expectedNames[i], name)
		}

		if obj == nil {
			continue
		}

		if obj.Name != expectedNames[i] {
			t.Errorf("names expected (%s) got (%s)", expectedNames[i], name)
		}

		if obj.Path != expectedPaths[i] {
			t.Errorf("paths expected (%s) got (%s)", expectedPaths[i], obj.Path)
		}

		if obj.ContentType != expectedContentTypes[i] {
			t.Errorf("contentType %d expected (%s) got (%s)",
				len(expectedContentTypes), expectedContentTypes[i], obj.ContentType)
		}
	}
}
