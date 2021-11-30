//go:build test
// +build test

package kinesis_stream

import (
	"context"
	"fmt"
	"reflect"
)

type mockKinesisStreamClientImpl struct {
	kinesisStreamClient *kinesis.Client
	rt                  string
	cs                  []string
}

func (m *mockKinesisStreamClientImpl) putRecord(_ context.Context, _ string, rt interface{}, csArray []string) error {
	if csArray == nil {
		return fmt.Errorf("Expect csItems to not be nil.")
	}
	if e, a := m.rt, rt; !reflect.DeepEqual(e, a) {
		return fmt.Errorf("expect %v, got %v", e, a)
	}
	if e, a := m.cs, csArray; !reflect.DeepEqual(e, a) {
		return fmt.Errorf("expect %v, got %v", e, a)
	}
	return nil
}
