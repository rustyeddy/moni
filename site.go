package main

import (
	"net/url"

	log "github.com/sirupsen/logrus"
)

type Sites map[url.URL]*Site

// Site is a holding structure consisting of a URL, pointer to the
// home page and a map of all child pages.
type Site struct {
	*url.URL
	*Page // Home page
	Pages map[url.URL]*Page
}

func NewSite(u *url.URL) *Site {
	return &Site{u, nil, make(map[url.URL]*Page)}
}

func AddSite(u *url.URL) (s *Site) {
	if s = NewSite(u); s != nil {
		sites[*u] = s
	}
	return s
}

// GetSite accepts either a string representing a URL or a parsed
// url.URL, either way that URL is scrubbed and matched against the
// ACL to determine if it is to be walked.  If so a site object is
// obtained and returned.
func GetSite(urlstr string) (s *Site) {
	var u *url.URL
	var err error

	if u, err = url.Parse(urlstr); err != nil {
		log.Errorf("converting url: %v", err)
		return
	}

	// If scheme is "" it will likely have the hostname as the path
	// assuming that localhost is meant, which in this case is not
	// true.  To over come this problem, I add an http scheme, then
	// parse the string generated from the URL with http as the scheme.
	if u.Scheme == "" {
		u.Scheme = "http"
		u, err = url.Parse(u.String())
		if err != nil {
			log.Errorf("Error with URL %+v", u)
			return
		}
	}

	var e bool
	if s, e = sites[*u]; !e {
		if s = NewSite(u); s != nil {
			sites[*u] = s
		}
	}

	return s
}

// setupSites takes a slice of strings that is assumed to be URLs. The
// strings are scrubbed, and converted to a url.URL if they are
// legit.  The URL is then matched against an access-list (ACL) to
// determine if the given URL will be walked or not.  If not, it is
// represented by a _blank-page_. If the URL is to be walked, it is
// added to the watchlist to be walked and scheduled for future walks.
func submitSites(slist []string) {
	for _, s := range slist {
		log.Infof("submitSites to walkQ <- %s", s)
		walkQ <- s
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
func saveSitesFile() (err error) {

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
	if s.HomePage() == nil {
		return nil
	}
	for l, _ := range s.Links {
		plist = append(plist, l)
	}
	return plist
}

// NewPage returns a page structure that will hold all our cool stuff
func (s *Site) NewPage(u url.URL) (p *Page) {
	p = NewPage(&u)
	p.Site = s
	s.Page = p     // Home Page
	s.Pages[u] = p // All Pages
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

// HomePage return the page associated with the sites root or baseURL.
// If a page currently does not exist we will create one.
func (s *Site) HomePage() (p *Page) {
	if s.Page == nil {
		s.Page = s.NewPage(*s.URL) // Make the home page
		s.Pages[*s.URL] = s.Page   // save home with rest of 'em
	}
	return s.Page
}

func (p *Page) String() string {
	return p.URL.String()
}
