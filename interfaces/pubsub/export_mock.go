//go:build test
// +build test

package pubsub

import (
	"cloud.google.com/go/pubsub"
	"context"
	"fmt"
	"reflect"
)

type mockPubsubClientImpl struct {
	pubsubClient *pubsub.Client
	cs           []string
}

func (_ *mockPubsubClientImpl) topicExists(ctx context.Context, topicID string) (bool, error) {
	return false, nil
}

func (_ *mockPubsubClientImpl) createTopic(ctx context.Context, topicID string) (*pubsub.Topic, error) {
	return nil, nil
}

func (_ *mockPubsubClientImpl) subscriptionExists(ctx context.Context, subscriptionID string) (bool, error) {
	return false, nil
}

func (_ *mockPubsubClientImpl) createSubscription(ctx context.Context, topicID string, subscriptionID string) (*pubsub.Subscription, error) {
	return nil, nil
}

func (m *mockPubsubClientImpl) publishMessage(_ context.Context, _ string, csArray []string) error {
	if csArray == nil {
		return fmt.Errorf("Expect csItems to not be nil.")
	}
	if e, a := m.cs, csArray; !reflect.DeepEqual(e, a) {
		return fmt.Errorf("expect %v, got %v", e, a)
	}
	return nil
}
