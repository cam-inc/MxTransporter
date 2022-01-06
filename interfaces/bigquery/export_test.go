//go:build test
// +build test

package bigquery

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"math"
	"testing"
	"time"
)

func Test_ExportToBigquery(t *testing.T) {
	csMap := primitive.M{
		"_id":               primitive.M{"_data": "00000"},
		"operationType":     "insert",
		"clusterTime":       primitive.Timestamp{00000, 0},
		"fullDocument":      primitive.M{"wwwww": "test full document"},
		"ns":                primitive.M{"xxxxx": "test ns"},
		"documentKey":       primitive.M{"yyyyy": "test document key"},
		"updateDescription": primitive.M{"zzzzz": "test update description"},
	}

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

	ctx := context.Background()

	tests := []struct {
		name string
		runner func(t *testing.T)
	}{
		{
			name: "Pass to put a record to bigquery.",
			runner: func(t *testing.T) {
				bqClientImpl := &mockBigqueryClientImpl{nil, testCsItems}
				mockBqImpl := BigqueryImpl{bqClientImpl}
				if err := mockBqImpl.ExportToBigquery(ctx, csMap); err != nil {
					t.Fatalf("Testing Error, ErrorMessage: %v", err)
				}
			},
		},
		{
			name: "Failed to marshal _id parameter of csMap.",
			runner: func(t *testing.T) {
				// Insert something that json.marchal fails
				csMap = primitive.M{
					"_id":               math.NaN(),
					"operationType":     "insert",
					"clusterTime":       primitive.Timestamp{00000, 0},
					"fullDocument":      primitive.M{"wwwww": "test full document"},
					"ns":                primitive.M{"xxxxx": "test ns"},
					"documentKey":       primitive.M{"yyyyy": "test document key"},
					"updateDescription": primitive.M{"zzzzz": "test update description"},
				}

				bqClientImpl := &mockBigqueryClientImpl{nil, nil}
				mockBqImpl := BigqueryImpl{bqClientImpl}
				if err := mockBqImpl.ExportToBigquery(ctx, csMap); err == nil {
					t.Fatalf("Not behaving as intended.")
				}
			},
		},
		{
			name: "Failed to marshal _id parameter of csMap.",
			runner: func(t *testing.T) {
				// Insert something that json.marchal fails
				csMap := primitive.M{
					"_id":               primitive.M{"_data": "00000"},
					"operationType":     "insert",
					"clusterTime":       primitive.Timestamp{00000, 0},
					"fullDocument":      math.NaN(),
					"ns":                primitive.M{"xxxxx": "test ns"},
					"documentKey":       primitive.M{"yyyyy": "test document key"},
					"updateDescription": primitive.M{"zzzzz": "test update description"},
				}

				bqClientImpl := &mockBigqueryClientImpl{nil, nil}
				mockBqImpl := BigqueryImpl{bqClientImpl}
				if err := mockBqImpl.ExportToBigquery(ctx, csMap); err == nil {
					t.Fatalf("Not behaving as intended.")
				}
			},
		},
		{
			name: "Failed to marshal _id parameter of csMap.",
			runner: func(t *testing.T) {
				// Insert something that json.marchal fails
				csMap := primitive.M{
					"_id":               primitive.M{"_data": "00000"},
					"operationType":     "insert",
					"clusterTime":       primitive.Timestamp{00000, 0},
					"fullDocument":      primitive.M{"wwwww": "test full document"},
					"ns":                math.NaN(),
					"documentKey":       primitive.M{"yyyyy": "test document key"},
					"updateDescription": primitive.M{"zzzzz": "test update description"},
				}

				bqClientImpl := &mockBigqueryClientImpl{nil, nil}
				mockBqImpl := BigqueryImpl{bqClientImpl}
				if err := mockBqImpl.ExportToBigquery(ctx, csMap); err == nil {
					t.Fatalf("Not behaving as intended.")
				}
			},
		},
		{
			name: "Failed to marshal _id parameter of csMap.",
			runner: func(t *testing.T) {
				// Insert something that json.marchal fails
				csMap := primitive.M{
					"_id":               primitive.M{"_data": "00000"},
					"operationType":     "insert",
					"clusterTime":       primitive.Timestamp{00000, 0},
					"fullDocument":      primitive.M{"wwwww": "test full document"},
					"ns":                primitive.M{"xxxxx": "test ns"},
					"documentKey":       math.NaN(),
					"updateDescription": primitive.M{"zzzzz": "test update description"},
				}

				bqClientImpl := &mockBigqueryClientImpl{nil, nil}
				mockBqImpl := BigqueryImpl{bqClientImpl}
				if err := mockBqImpl.ExportToBigquery(ctx, csMap); err == nil {
					t.Fatalf("Not behaving as intended.")
				}
			},
		},
		{
			name: "Failed to marshal _id parameter of csMap.",
			runner: func(t *testing.T) {
				// Insert something that json.marchal fails
				csMap := primitive.M{
					"_id":               primitive.M{"_data": "00000"},
					"operationType":     "insert",
					"clusterTime":       primitive.Timestamp{00000, 0},
					"fullDocument":      primitive.M{"wwwww": "test full document"},
					"ns":                primitive.M{"xxxxx": "test ns"},
					"documentKey":       primitive.M{"yyyyy": "test document key"},
					"updateDescription": math.NaN(),
				}

				bqClientImpl := &mockBigqueryClientImpl{nil, nil}
				mockBqImpl := BigqueryImpl{bqClientImpl}
				if err := mockBqImpl.ExportToBigquery(ctx, csMap); err == nil {
					t.Fatalf("Not behaving as intended.")
				}
			},
		},
	}

	for _, v := range tests {
		t.Run(v.name, v.runner)
	}
}
