package mongodb

import (
	"os"
)

type Mongo struct {
	MongoDbConnectionUrl string
	MongoDbDatabase      string
	MongoDbCollection    string
}

func MongoConfig() Mongo {
	var mCfg Mongo
	mCfg.MongoDbConnectionUrl = os.Getenv("MONGODB_HOST")
	mCfg.MongoDbDatabase = os.Getenv("MONGODB_DATABASE")
	mCfg.MongoDbCollection = os.Getenv("MONGODB_COLLECTION")
	return mCfg
}
