package moni

import (
	"context"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
)

type MDB struct {
	*mongo.Database
	*mongo.Client
}

var (
	mdb MDB = MDB{}
	err error
)

// StartDatabase will start the database
func init() {
	var err error
	cli, err := mongo.Connect(context.Background(), "mongodb://localhost:27017", nil)
	IfErrorFatal(err, "mongo connect")
	mdb.Database = cli.Database("monty")
}

func DropDatabase(name string) {
	_, err = mdb.RunCommand(
		context.Background(),
		bson.NewDocument(bson.EC.Int32("dropDatabase", 1)),
	)
}

func FetchStorage() (si *StorageInfo) {
	panic("must implement")
	return si
}

func (mdb *MDB) connect() {
	cli, err := mongo.Connect(context.Background(), "mongodb://localhost:27017", nil)
	IfErrorFatal(err, "mongo connect")
	mdb.Database = cli.Database("monty")
}

func (mdb *MDB) UseStore(name string) *mongo.Collection {
	coll := mdb.Collection(name)
	return coll
}
