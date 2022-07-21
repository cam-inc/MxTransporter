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
		if err := os.Setenv("PUBSUB_TOPIC_NAME", "xxx"); err != nil {
			t.Fatalf("Failed to set file PUBSUB_TOPIC_NAME environment variables.")
		}
		if err := os.Setenv("PUBSUB_ORDERING_BY", "yyy"); err != nil {
			t.Fatalf("Failed to set file PUBSUB_ORDERING_BY environment variables.")
		}
		psCfg := PubSubConfig()
		want := PubSub{
			TopicName:  "xxx",
			OrderingBy: "yyy",
		}
		if !reflect.DeepEqual(want, psCfg) {
			t.Fatalf("Environment variable PUBSUB_* is not acquired correctly. want: %v, got: %v", want, psCfg)
		}
	})
}
