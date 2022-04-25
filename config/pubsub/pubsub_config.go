package pubsub

import (
	"github.com/cam-inc/mxtransporter/config/constant"
	"os"
)

type PubSub struct {
	TopicName string
}

func PubSubConfig() PubSub {
	var psCfg PubSub
	psCfg.TopicName = os.Getenv(constant.PUBSUB_TOPIC_NAME)
	return psCfg
}
