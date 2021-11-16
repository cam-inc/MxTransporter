package pubsub

import (
	"cloud.google.com/go/pubsub"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"reflect"
	"strconv"
	"testing"
	"time"
)

var csMap = primitive.M{
	"_id":               primitive.M{"_data": "00000"},
	"operationType":     "insert",
	"clusterTime":       primitive.Timestamp{00000, 0},
	"fullDocument":      primitive.M{"wwwww": "test full document"},
	"ns":                primitive.M{"xxxxx": "test ns"},
	"documentKey":       primitive.M{"yyyyy": "test document key"},
	"updateDescription": primitive.M{"zzzzz": "test update description"},
}

//type mockPubsubFuncs struct{}
type mockPubsubClientImpl struct {
	PubsubClient
	fakePubsubTopic func(ctx context.Context, topicID string, psClient *pubsub.Client) error
	fakePubsubSubscription func(ctx context.Context, topicID string, subscriptionID string, psClient *pubsub.Client) error
	fakePublishMessage func(ctx context.Context, topicID string, csArray []string, psClient *pubsub.Client) error
}

func (m *mockPubsubClientImpl) PubsubTopic(ctx context.Context, topicID string, psClient *pubsub.Client) error{
	return m.fakePubsubTopic(ctx, topicID, psClient)
}

func (m *mockPubsubClientImpl) PubsubSubscription(ctx context.Context, topicID string, subscriptionID string, psClient *pubsub.Client) error{
	return m.fakePubsubSubscription(ctx, topicID, subscriptionID, psClient)
}

func (m *mockPubsubClientImpl) PublishMessage(ctx context.Context, topicID string, csArray []string, psClient *pubsub.Client) error{
	return m.fakePublishMessage(ctx, topicID, csArray, psClient)
}

func Test_ExportToPubSub(t *testing.T) {

	cases := []struct {
		cs       primitive.M
		client *pubsub.Client
		function PubsubClient
	}{
		{
			cs:       csMap,
			client: nil,
			function: &mockPubsubClientImpl{
				fakePubsubTopic: func(_ context.Context, _ string, _ *pubsub.Client) error {
					return nil
				},
				fakePubsubSubscription: func(_ context.Context, _ string, _ string, _ *pubsub.Client) error {
					return nil
				},
				fakePublishMessage: func(_ context.Context, _ string, csArray []string, _ *pubsub.Client) error {
					testCsArray := []string{
						`{"_data":"00000"}`,
						"insert",
						time.Unix(int64(csMap["clusterTime"].(primitive.Timestamp).T), 0).Format("2006-01-02 15:04:05"),
						`{"wwwww":"test full document"}`,
						`{"xxxxx":"test ns"}`,
						`{"yyyyy":"test document key"}`,
						`{"zzzzz":"test update description"}`,
					}

					if csArray == nil {
						return fmt.Errorf("expect csItems to not be nil")
					}
					if e, a := testCsArray, csArray; !reflect.DeepEqual(e, a) {
						return fmt.Errorf("expect %v, got %v", e, a)
					}
					return nil
				},
			},
		},
	}

	for i, tt := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			ctx := context.TODO()
			if err := NewPubsubClient(tt.function).ExportToPubSub(ctx, tt.cs, tt.client); err != nil {
				t.Fatalf("expect no error, got %v", err)
			}
		})
	}
}
