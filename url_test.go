package moni

import "testing"

func TestNormalizeURL(t *testing.T) {
	var tsts = []struct {
		host     string
		expected string
	}{
		{"amazon.com", "http://amazon.com"},
		{"//clowdops.net:4040", "http://clowdops.net:4040"},
		{"http://rustyeddy.com", "http://rustyeddy.com"},
		{"//example.com:300", "http://example.com:300"},
		{"//john.bozo.nono.com:3000", "http://john.bozo.nono.com:3000"},
		{"tel:phonenumber", ""},
		{"", ""},
	}

	for _, tst := range tsts {
		ustr, err := NormalizeURL(tst.host)
		if err != nil && tst.expected != "" {
			// Here if we have an error and did not expect one
			t.Errorf("Normalize URL failed %v", err)
		}
		if ustr != tst.expected {
			t.Errorf("Normalize URL (%s) expected (%s) got (%s)", tst.host, tst.expected, ustr)
		}
	}
}
