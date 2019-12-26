package main

import (
	"net/url"
	"testing"
)

func TestNewPage(t *testing.T) {
	u, _ := url.Parse("example.com")
	p := NewPage(u)
	if p == nil {
		t.Error("expected to create a page but failed")
	}

	if len(p.Links) > 0 {
		t.Errorf("TestNewPage expected (0) links got (%d)", len(p.Links))
	}
}
