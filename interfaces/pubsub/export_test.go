package pubsub

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"reflect"
	"strconv"
	"testing"
	"time"
)

var csMap = primitive.M{
	"_id": primitive.M{"_data": "00000"},
	"operationType": "insert",
	"clusterTime": primitive.Timestamp{00000, 0},
	"fullDocument": primitive.M{"xxxxx": "xxxxx"},
	"ns": primitive.M{"xxxxx": "xxxxx"},
	"documentKey": primitive.M{"xxxxx": "xxxxx"},
	"updateDescription": primitive.M{"xxxxx": "xxxxx"},
}

type MockPubsubFuncs struct {}

func (m *MockPubsubFuncs) PubsubTopic(ctx context.Context, topicId string) error {
	_ = func (t *testing.T) error {
		t.Helper()
		if topicId == "" {
			t.Fatal("expect topicId to not be nil")
		}
		return nil
	}
	return nil
}

func (m *MockPubsubFuncs) PubsubSubscription(ctx context.Context, topicId string, subscriptionId string) error {
	_ = func (t *testing.T) error {
		t.Helper()
		if topicId == "" {
			t.Fatal("expect topicId to not be nil")
		}
		if subscriptionId == "" {
			t.Fatal("expect subscriptionId to not be nil")
		}
		return nil
	}
	return nil
}

func (m *MockPubsubFuncs) PublishMessage(ctx context.Context, topicId string, csArray []string) error {
	_ = func (t *testing.T) error {
		t.Helper()

		testCsArray := []string{
			"{\"_data\":\"00000\"}",
			"insert",
			time.Unix(int64(csMap["clusterTime"].(primitive.Timestamp).T), 0).Format("2006-01-02 15:04:05"),
			"{\"xxxxx\":\"xxxxx\"}",
			"{\"xxxxx\":\"xxxxx\"}",
			"{\"xxxxx\":\"xxxxx\"}",
			"{\"xxxxx\":\"xxxxx\"}",
		}

		if topicId == "" {
			t.Fatal("expect topicId to not be nil")
		}
		if csArray == nil {
			t.Fatal("expect csItems to not be nil")
		}
		if e, a := testCsArray, csArray; !reflect.DeepEqual(e, a) {
			t.Errorf("expect %v, got %v", e, a)
		}
		return nil
	}
	return nil
}

func Test_ExportToPubSub(t *testing.T) {

	cases := []struct {
		cs primitive.M
	}{
		{
			cs: csMap,
		},
	}

	function := &MockPubsubFuncs{}

	for i, tt := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			ctx := context.TODO()
			err := ExportToPubSub(ctx, tt.cs, function)
			if err != nil {
				t.Fatalf("expect no error, got %v", err)
			}
		})
	}
}

