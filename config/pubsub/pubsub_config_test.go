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
		mDatabase := "xxx"
		mCollection := "yyy"
		if err := os.Setenv("MONGODB_DATABASE", mDatabase); err != nil {
			t.Fatalf("Failed to set file MONGODB_DATABASE environment variables.")
		}
		if err := os.Setenv("MONGODB_COLLECTION", mCollection); err != nil {
			t.Fatalf("Failed to set file MONGODB_COLLECTION environment variables.")
		}

		psCfg := PubSubConfig()
		if e, a := psCfg.MongoDbDatabase, mDatabase; !reflect.DeepEqual(e, a) {
			t.Fatal("Environment variable MONGODB_DATABASE is not acquired correctly.")
		}
		if e, a := psCfg.MongoDbCollection, mCollection; !reflect.DeepEqual(e, a) {
			t.Fatal("Environment variable MONGODB_COLLECTION is not acquired correctly.")
		}
	})
}