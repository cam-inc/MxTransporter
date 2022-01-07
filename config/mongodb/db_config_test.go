//go:build test
// +build test

package mongodb

import (
	"os"
	"reflect"
	"testing"
)

func Test_MongoConfig(t *testing.T) {
	t.Run("Check to call the set environment variable.", func(t *testing.T) {
		mConnectionUrl := "mongodb+srv://xxx:yyy@user.www.mongodb.net"
		mDatabase := "xxx"
		mCollection := "yyy"
		if err := os.Setenv("MONGODB_HOST", mConnectionUrl); err != nil {
			t.Fatalf("Failed to set file MONGODB_HOST environment variables.")
		}
		if err := os.Setenv("MONGODB_DATABASE", mDatabase); err != nil {
			t.Fatalf("Failed to set file MONGODB_DATABASE environment variables.")
		}
		if err := os.Setenv("MONGODB_COLLECTION", mCollection); err != nil {
			t.Fatalf("Failed to set file MONGODB_COLLECTION environment variables.")
		}

		mCfg := MongoConfig()
		if e, a := mCfg.MongoDbConnectionUrl, mConnectionUrl; !reflect.DeepEqual(e, a) {
			t.Fatal("Environment variable MONGODB_HOST is not acquired correctly.")
		}
		if e, a := mCfg.MongoDbDatabase, mDatabase; !reflect.DeepEqual(e, a) {
			t.Fatal("Environment variable MONGODB_DATABASE is not acquired correctly.")
		}
		if e, a := mCfg.MongoDbCollection, mCollection; !reflect.DeepEqual(e, a) {
			t.Fatal("Environment variable MONGODB_COLLECTION is not acquired correctly.")
		}
	})
}