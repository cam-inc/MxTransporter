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
	var mongoCfg Mongo
	mongoCfg.MongoDbConnectionUrl = os.Getenv("MONGODB_HOST")
	mongoCfg.MongoDbDatabase = os.Getenv("MONGODB_DATABASE")
	mongoCfg.MongoDbCollection = os.Getenv("MONGODB_COLLECTION")
	return mongoCfg
}
