package pubsub

import (
	"cloud.google.com/go/pubsub"
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strconv"
	"testing"
)

type mockPubsubTopic func(id string) *pubsub.Topic
type mockPubsubCreatetopic func(ctx context.Context, topicID string) (*pubsub.Topic, error)
type mockPubsubSubscription func(id string) *pubsub.Subscription
type mockCreateSubscription func(ctx context.Context, id string, cfg pubsub.SubscriptionConfig) (*pubsub.Subscription, error)

func (m mockPubsubTopic) Topic(id string) *pubsub.Topic {
	return m(id)
}

func (m mockPubsubCreatetopic) CreateTopic(ctx context.Context, topicID string) (*pubsub.Topic, error) {
	return m(ctx, topicID)
}

func (m mockPubsubSubscription) Subscription(id string) *pubsub.Subscription {
	return m(id)
}

func (m mockCreateSubscription) CreateSubscription(ctx context.Context, id string, cfg pubsub.SubscriptionConfig) (*pubsub.Subscription, error) {
	return m(ctx, id, cfg)
}

func Test_ExportToPubSub(t *testing.T) {
	csMap := primitive.M{
		"_id": primitive.M{"_data": "00000"},
		"operationType": "insert",
		"clusterTime": primitive.Timestamp{00000, 0},
		"fullDocument": primitive.M{"xxxxx": "xxxxx"},
		"ns": primitive.M{"xxxxx": "xxxxx"},
		"documentKey": primitive.M{"xxxxx": "xxxxx"},
		"updateDescription": primitive.M{"xxxxx": "xxxxx"},
	}

	cases := []struct {
		client func(t *testing.T) pubsubClient
		cs primitive.M
	}{
		{
			client: func(t *testing.T) pubsubClient {
				return {
					mockPubsubTopic(func(id string) *pubsub.Topic {
						t.Helper()
						return nil
					}),
					mockPubsubCreatetopic(func(ctx context.Context, topicID string) (*pubsub.Topic, error) {
						t.Helper()
						return nil, nil
					}),
					mockPubsubSubscription(func(id string) *pubsub.Subscription {
					}),
				}
				return nil
			},
			cs: csMap,
		},
	}

	for i, tt := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			ctx := context.TODO()
			err := ExportToPubSub(ctx, tt.cs, tt.client(t))
			if err != nil {
				t.Fatalf("expect no error, got %v", err)
			}
		})
	}
}

