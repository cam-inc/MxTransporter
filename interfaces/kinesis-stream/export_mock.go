//go:build test
// +build test

package kinesis_stream

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/kinesis"
	"reflect"
)

type mockKinesisStreamClientImpl struct {
	kinesisStreamClient *kinesis.Client
	rt                  string
	cs                  []string
}

type mockKinesisStreamClientImplError struct {
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

func (m *mockKinesisStreamClientImplError) putRecord(_ context.Context, _ string, _ interface{}, _ []string) error {
	return fmt.Errorf("Expected errors for error handling.")
}