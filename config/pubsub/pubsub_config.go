package pubsub

import (
	"os"
)

type PubSub struct {
	MongoDbDatabase   string
	MongoDbCollection string
}

func PubSubConfig() PubSub {
	var pubSubConfig PubSub
	pubSubConfig.MongoDbDatabase = os.Getenv("MONGODB_DATABASE")
	pubSubConfig.MongoDbCollection = os.Getenv("MONGODB_COLLECTION")
	return pubSubConfig
}
