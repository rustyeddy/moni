package main

import "testing"

func TestACL(t *testing.T) {
	inputs := []struct {
		inurl string
		ex    bool
	}{
		{"google.com", false},
		{"example.com", false},
		{"gumsole.com", true},
	}

	for _, b := range inputs {
		allow := acl.Allow(b.inurl)
		if allow != b.ex {
			t.Errorf("test inputs for %s expected (%v) got (%v)", b.inurl, allow, b.ex)
		}
	}
}

func TestAddACL(t *testing.T) {
	inputs := []struct {
		inurl string
		ex    bool
	}{
		{"example.com", false},
		{"foo.bar", false},
		{"bar.baz", false},
	}

	for _, b := range inputs {
		allow := acl.Allow(b.inurl)
		if allow != b.ex {
			t.Errorf("test inputs for %s expected (%v) got (%v)", b.inurl, allow, b.ex)
		}
	}

	acl.Add("example.com", true)
	acl.Add("foo.bar", true)

	inputs[0].ex = true
	inputs[1].ex = true

	for _, b := range inputs {
		allow := acl.Allow(b.inurl)
		if allow != b.ex {
			t.Errorf("test inputs for %s expected (%v) got (%v)", b.inurl, allow, b.ex)
		}
	}

}
