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
	var mongoConfig Mongo
	mongoConfig.MongoDbConnectionUrl = os.Getenv("MONGODB_HOST")
	mongoConfig.MongoDbDatabase = os.Getenv("MONGODB_DATABASE")
	mongoConfig.MongoDbCollection = os.Getenv("MONGODB_COLLECTION")
	return mongoConfig
}
