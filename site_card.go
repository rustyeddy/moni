package moni

// WebUI
// ====================================================================

// SitesCard provides a card with list of the sites we are managing
type SitesCard struct {
	*Card
	Sitemap
}

func (app *App) NewSitesCard() (c *SitesCard) {
	c = &SitesCard{
		Card:    NewCard("Sites"),
		Sitemap: sites,
	}
	return c
}
