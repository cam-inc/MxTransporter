//go:build test
// +build test

package bigquery

import (
	"cloud.google.com/go/bigquery"
	"context"
	"fmt"
	"reflect"
)

type mockBigqueryClientImpl struct {
	bqClient *bigquery.Client
	csItems  []ChangeStreamTableSchema
}

type mockBigqueryClientImplError struct {
	bqClient *bigquery.Client
	csItems  []ChangeStreamTableSchema
}

func (m *mockBigqueryClientImpl) putRecord(_ context.Context, _ string, _ string, csItems []ChangeStreamTableSchema) error {
	if csItems == nil {
		return fmt.Errorf("Expect csItems to not be nil.")
	}

	if e, a := m.csItems, csItems; !reflect.DeepEqual(e, a) {
		return fmt.Errorf("expect %v, got %v", e, a)
	}
	return nil
}

func (m *mockBigqueryClientImplError) putRecord(_ context.Context, _ string, _ string, _ []ChangeStreamTableSchema) error {
	return fmt.Errorf("Expected errors for error handling.")
}
