package pubsub

import (
	"github.com/cam-inc/mxtransporter/config/constant"
	"os"
)

type PubSub struct {
	MongoDbDatabase   string
	MongoDbCollection string
}

func PubSubConfig() PubSub {
	var psCfg PubSub
	psCfg.MongoDbDatabase = os.Getenv(constant.MONGODB_DATABASE)
	psCfg.MongoDbCollection = os.Getenv(constant.MONGODB_COLLECTION)
	return psCfg
}
