package moni

import "testing"

func TestNewPage(t *testing.T) {
	p := NewPage("example.com")
	if p == nil {
		t.Error("creating page")
	}

	url := p.URL
	if url != "example" {
		t.Errorf("p.URL expected (example.com) got (%s) ", url)
	}
	if p.Links == nil || len(p.Links) != 0 {
		t.Errorf("expected p.Links != nil && len == 10 got %+v", p.Links)
	}
	if p.Ignored == nil || len(p.Ignored) != 0 {
		t.Errorf("expected p.Ignored != nil && len == 10 got %+v", p.Ignored)
	}

	if p.CrawlReady != true {
		t.Errorf("p.CrawlReady expected (true) got (%s) ", "false")
	}

	if len(pages) != 1 {
		t.Errorf("paged expected (%d) got (%d) ", 1, len(pages))
	}

	if _, ex := pages[url]; !ex {
		t.Errorf("expected %s to be an index paged got (not there) ", url)
	}

	// Now get page
	if p := GetPage(url); p == nil {
		t.Errorf("failed to get page for %s ", url)
	}
}
