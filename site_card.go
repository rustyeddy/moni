package moni

// WebUI
// ====================================================================

// SitesCard provides a card with list of the sites we are managing
type SitesCard struct {
	*Card
	Sitemap
}

func NewSitesCard() (c *SitesCard) {
	c = &SitesCard{
		Card:    NewCard("Sites"),
		Sitemap: Sites,
	}
	return c
}
