//go:build test
// +build test

package pubsub

import (
	"context"
	"fmt"
	"reflect"
)

func (_ *mockPubsubClientImpl) pubsubTopic(_ context.Context, _ string) error {
	return nil
}

func (_ *mockPubsubClientImpl) pubsubSubscription(_ context.Context, _ string, _ string) error {
	return nil
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
