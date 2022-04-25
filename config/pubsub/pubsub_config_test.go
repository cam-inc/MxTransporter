//go:build test
// +build test

package pubsub

import (
	"os"
	"reflect"
	"testing"
)

func Test_PubSubConfig(t *testing.T) {
	t.Run("Check to call the set environment variable.", func(t *testing.T) {
		mTopicID := "xxx"
		if err := os.Setenv("PUBSUB_TOPIC_NAME", mTopicID); err != nil {
			t.Fatalf("Failed to set file PUBSUB_TOPIC_NAME environment variables.")
		}

		psCfg := PubSubConfig()
		if e, a := psCfg.MongoDbDatabase, mTopicID; !reflect.DeepEqual(e, a) {
			t.Fatal("Environment variable MONGODB_DATABASE is not acquired correctly.")
		}
	})
}
