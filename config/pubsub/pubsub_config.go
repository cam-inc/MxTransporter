package pubsub

import (
	"os"

	"github.com/cam-inc/mxtransporter/config/constant"
)

type PubSub struct {
	TopicName  string
	OrderingBy string
}

func PubSubConfig() PubSub {
	var psCfg PubSub
	psCfg.TopicName = os.Getenv(constant.PUBSUB_TOPIC_NAME)
	psCfg.OrderingBy = os.Getenv(constant.PUBSUB_ORDERING_BY)
	return psCfg
}
