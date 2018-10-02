package inv

type SiteMap map[string]*Site

struct Site struct {
	Baseurl string
	PageMap
}
