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

func (*mockPubsubClientImpl) topicExists(ctx context.Context, topicID string) (bool, error) {
	return false, nil
}

func (*mockPubsubClientImpl) createTopic(ctx context.Context, topicID string) (*pubsub.Topic, error) {
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
