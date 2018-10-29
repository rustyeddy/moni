package moni

import (
	"context"
	"time"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
	log "github.com/sirupsen/logrus"
)

// Site is basically a website wich includes API interfaces
type Site struct {
	ID      int64  `bson:"_id" json:"id"`
	URL     string `bson:"url" json:"url"`
	IP      string `bson:"ip" json:"ip"`
	Health  bool   `bson:"health" json:"health"`
	Pagemap `bson:"pagemap" json:"pagemap"`

	// Crawl job info
	LastCrawled time.Time `bson:"last_crawled" json:"last_crawled"`

	nextCrawl  time.Time // Ignore this in
	crawlState int
	*time.Timer
	*log.Entry // Ignore
}

// SiteManager
// ====================================================================
func siteCollection() (col *mongo.Collection) {
	col = mdb.Collection("sites")
	IfNilFatal(col, "GetSitesCollection")
	return col
}

func FetchSites() (sites []*Site) {
	col := siteCollection()
	cur, err := col.Find(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}
	defer cur.Close(context.Background())
	for cur.Next(context.Background()) {
		elem := bson.NewDocument()
		if err := cur.Decode(elem); err != nil {
			log.Fatal(err)
		}
		log.Fatalf("cur: %+v\n", cur)
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}
	return sites
}

func FetchSite(url string) *Site {
	//site := bson.NewDocument()
	filter := bson.NewDocument(bson.EC.String("URL", url))
	col := siteCollection()

	var site Site
	err := col.FindOne(context.Background(), filter).Decode(&site)
	IfErrorFatal(err)

	log.Fatalf("result %+v", site)
	return &site
}

// map[string]string{"hello": "world"})
func StoreSite(s *Site) (id int64) {
	col := siteCollection()
	if res, err := col.InsertOne(context.Background(), s); err != nil {
		IfErrorFatal(err)
	} else {
		s.ID = res.InsertedID.(int64)
	}
	return s.ID
}

func StoreManySites(s []*Site) (ids []int64) {
	panic("Todo need to implement")
	return ids
}

func DeleteSite(url string) {
	panic("todo")
}

// AddNewSite will create a New Site from the url, including
// verify and sanitize the url and so on.
func AddUrl(url string) {

	// Get the site collection to be added to
	col := siteCollection()
	IfNilError(col, "AddSite")

	// Schedule a new crawl
	// Store the site
	log.Infof("Added new site %s ~ calling StoreSites()", url)

	Crawler.UrlQ <- url
}

// ====================================================================

// ScheduleCrawl
func (s *Site) ScheduleCrawl() {
	timer := time.AfterFunc(time.Minute*5, func() {
		s.crawlState = CrawlReady
	})
	defer timer.Stop()
}
