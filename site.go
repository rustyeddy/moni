package main

import (
	"net/url"

	log "github.com/sirupsen/logrus"
)

type Sites map[url.URL]*Site

type Site struct {
	*url.URL
	Pages map[url.URL]*Page
}

func NewSite(u *url.URL) *Site {
	return &Site{u, make(map[url.URL]*Page)}
}

func AddSite(u *url.URL) (s *Site) {
	if s = NewSite(u); s != nil {
		sites[*u] = s
	}
	return s
}

func GetSite(uval interface{}) (s *Site) {
	var u *url.URL
	var err error

	// Make sure we have a url
	switch t := uval.(type) {
	case *url.URL:
		u = uval.(*url.URL)
	case string:
		urlstr := uval.(string)
		if u, err = url.Parse(urlstr); err != nil {
			log.Errorf("converting url: %v", err)
			return
		}
	default:
		log.Errorf("unknown url type: %s", t)
	}

	var e bool
	if s, e = sites[*u]; !e {
		if s = NewSite(u); s != nil {
			sites[*u] = s
		}
	}
	return s
}

func setupSites(wQ chan *Page, plist []string) {
	slist := readSitesFile()
	for _, urlstr := range append(slist, plist...) {
		if u := scrubURL(urlstr); u != nil {
			if site := GetSite(u); site != nil {
				if page := site.HomePage(); page != nil {
					wQ <- page
				}
			}
		}
	}
}

func readSitesFile() []string {
	sitelist := make([]string, 1)
	err := storage.ReadObject("sites.json", &sitelist)
	if err != nil {
		log.Errorf("Storage failed to read sites.json: %v", err)
		return nil
	}
	return sitelist
}

// SaveSites saves the sites structure.
func SaveSitesFile() (err error) {

	var sitelist []string
	sitelist = make([]string, 1)

	for _, s := range sites {
		urlstr := s.URL.String()
		sitelist = append(sitelist, urlstr)
	}
	if sitelist != nil && len(sitelist) > 0 {
		if err = storage.Save("sites.json", &sitelist); err != nil {
			log.Errorf("Storage Save failed for sites.json %v", err)
		}
	}
	return err
}

func (s *Site) PageList() (plist []string) {
	for _, p := range s.Pages {
		plist = append(plist, p.URL.String())
	}
	return plist
}

func (s *Site) HomePage() (p *Page) {
	p = s.GetPage(*s.URL)
	return p
}

// NewPage returns a page structure that will hold all our cool stuff
func (s *Site) NewPage(u url.URL) (p *Page) {
	p = &Page{
		Site:  s,
		URL:   u,
		Links: make(map[string][]string),
	}
	log.Infof("New Page: %+v for a total of %d\n", u, len(s.Pages))
	s.Pages[u] = p
	return p
}

// GetPage will return the page if it exists, or create otherwise.
func (s *Site) GetPage(u url.URL) (p *Page) {
	var ex bool
	if p, ex = s.Pages[u]; ex {
		return p
	}
	p = s.NewPage(u)
	return p
}

func (p *Page) String() string {
	return p.URL.String()
}
