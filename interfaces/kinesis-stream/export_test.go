//go:build test
// +build test

package kinesis_stream

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"math"
	"testing"
	"time"
)

func Test_ExportToKinesisStream(t *testing.T) {
	testRt := "00000"

	csMap := primitive.M{
		"_id":               primitive.M{"_data": "00000"},
		"operationType":     "insert",
		"clusterTime":       primitive.Timestamp{00000, 0},
		"fullDocument":      primitive.M{"wwwww": "test full document"},
		"ns":                primitive.M{"xxxxx": "test ns"},
		"documentKey":       primitive.M{"yyyyy": "test document key"},
		"updateDescription": primitive.M{"zzzzz": "test update description"},
	}

	testCsArray := []string{
		`{"_data":"00000"}`,
		"insert",
		time.Unix(int64(csMap["clusterTime"].(primitive.Timestamp).T), 0).Format("2006-01-02 15:04:05"),
		`{"wwwww":"test full document"}`,
		`{"xxxxx":"test ns"}`,
		`{"yyyyy":"test document key"}`,
		`{"zzzzz":"test update description"}`,
	}

	ctx := context.Background()


	tests := []struct {
		name string
		runner func(t *testing.T)
	}{
		{
			name: "Pass to put a record to kinesis data streams.",
			runner: func(t *testing.T) {
				ksClientImpl := &mockKinesisStreamClientImpl{nil, testRt, testCsArray}
				mockKsImpl := KinesisStreamImpl{ksClientImpl}
				if err := mockKsImpl.ExportToKinesisStream(ctx, csMap); err != nil {
					t.Fatalf("Testing Error, ErrorMessage: %v", err)
				}
			},
		},
		{
			name: "Failed to put a record to kinesis data streams.",
			runner: func(t *testing.T) {
				ksClientImpl := &mockKinesisStreamClientImplError{nil, "", nil}
				mockKsImpl := KinesisStreamImpl{ksClientImpl}
				if err := mockKsImpl.ExportToKinesisStream(ctx, csMap); err == nil {
					t.Fatalf("Not behaving as intended.")
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

				ksClientImpl := &mockKinesisStreamClientImpl{nil, "", nil}
				mockKsImpl := KinesisStreamImpl{ksClientImpl}
				if err := mockKsImpl.ExportToKinesisStream(ctx, csMap); err == nil {
					t.Fatalf("Not behaving as intended.")
				}
			},
		},
		{
			name: "Failed to marshal fullDocument parameter of csMap.",
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

				ksClientImpl := &mockKinesisStreamClientImpl{nil, "", nil}
				mockKsImpl := KinesisStreamImpl{ksClientImpl}
				if err := mockKsImpl.ExportToKinesisStream(ctx, csMap); err == nil {
					t.Fatalf("Not behaving as intended.")
				}
			},
		},
		{
			name: "Failed to marshal ns parameter of csMap.",
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

				ksClientImpl := &mockKinesisStreamClientImpl{nil, "", nil}
				mockKsImpl := KinesisStreamImpl{ksClientImpl}
				if err := mockKsImpl.ExportToKinesisStream(ctx, csMap); err == nil {
					t.Fatalf("Not behaving as intended.")
				}
			},
		},
		{
			name: "Failed to marshal documentKey parameter of csMap.",
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

				ksClientImpl := &mockKinesisStreamClientImpl{nil, "", nil}
				mockKsImpl := KinesisStreamImpl{ksClientImpl}
				if err := mockKsImpl.ExportToKinesisStream(ctx, csMap); err == nil {
					t.Fatalf("Not behaving as intended.")
				}
			},
		},
		{
			name: "Failed to marshal updateDescription parameter of csMap.",
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

				ksClientImpl := &mockKinesisStreamClientImpl{nil, "", nil}
				mockKsImpl := KinesisStreamImpl{ksClientImpl}
				if err := mockKsImpl.ExportToKinesisStream(ctx, csMap); err == nil {
					t.Fatalf("Not behaving as intended.")
				}
			},
		},
	}

	for _, v := range tests {
		t.Run(v.name, v.runner)
	}
}
