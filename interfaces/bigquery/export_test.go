//go:build test
// +build test

package bigquery

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

	t.Run("Test if the format of change streams works.", func(t *testing.T) {
		ctx := context.TODO()
		bqClientImpl := &mockBigqueryClientImpl{nil, testCsItems}
		mockBqImpl := BigqueryImpl{bqClientImpl}
		if err := mockBqImpl.ExportToBigquery(ctx, csMap); err != nil {
			t.Fatalf("Testing Error, ErrorMessage: %v", err)
		}
	})
}
