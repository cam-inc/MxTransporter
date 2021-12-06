package pubsub

import (
	"os"
)

type PubSub struct {
	MongoDbDatabase   string
	MongoDbCollection string
}

func PubSubConfig() PubSub {
	var psCfg PubSub
	psCfg.MongoDbDatabase = os.Getenv("MONGODB_DATABASE")
	psCfg.MongoDbCollection = os.Getenv("MONGODB_COLLECTION")
	return psCfg
}
