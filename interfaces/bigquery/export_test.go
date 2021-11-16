package bigquery

import (
	"cloud.google.com/go/bigquery"
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

type mockBigqueryClientImpl struct {
	BigqueryClient
	fakePutRecord func(ctx context.Context, dataset string, table string, csItems []ChangeStreamTableSchema, bqClient *bigquery.Client) error
}

func (m *mockBigqueryClientImpl) PutRecord(ctx context.Context, dataset string, table string, csItems []ChangeStreamTableSchema, bqClient *bigquery.Client) error {
	return m.fakePutRecord(ctx, dataset, table, csItems, bqClient)
}

func Test_ExportToBigquery(t *testing.T) {
	testCsItems := []ChangeStreamTableSchema{
		{
			ID:                `{"_data":"00000"}`,
			OperationType:     "insert",
			ClusterTime:       time.Unix(int64(csMap["clusterTime"].(primitive.Timestamp).T), 0),
			FullDocument:      `{"wwwww":"test full document"}`,
			Ns:                `{"xxxxx":"test ns"}`,
			DocumentKey:       `{"yyyyy":"test document key"}`,
			UpdateDescription: `{"zzzzz":"test update description"}`,
		},
	}

	cases := []struct {
		cs       primitive.M
		client   *bigquery.Client
		function BigqueryClient
	}{
		{
			cs:     csMap,
			client: nil,
			function: &mockBigqueryClientImpl{
				fakePutRecord: func(_ context.Context, _ string, _ string, csItems []ChangeStreamTableSchema, bqClient *bigquery.Client) error {
					if csItems == nil {
						return fmt.Errorf("expect csItems to not be nil")
					}
					if e, a := testCsItems, csItems; !reflect.DeepEqual(e, a) {
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
			if err := NewBigqueryClient(tt.function).ExportToBigquery(ctx, tt.cs, tt.client); err != nil {
				t.Fatalf("expect no error, got %v", err)
			}
		})
	}
}
