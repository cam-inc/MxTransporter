//go:build test
// +build test

package pubsub

import (
	"context"
	"math"
	"testing"
	"time"

	"github.com/cam-inc/mxtransporter/config"
	"github.com/cam-inc/mxtransporter/pkg/logger"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

func Test_ExportToPubSub(t *testing.T) {
	var l *zap.SugaredLogger
	logConfig := config.LogConfig()
	l = logger.New(logConfig)

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
		name   string
		runner func(t *testing.T)
	}{
		{
			name: "Pass to publish a message to pubsub.",
			runner: func(t *testing.T) {
				psClientImpl := &mockPubsubClientImpl{pubsubClient: nil, cs: testCsArray}
				mockPsImpl := PubsubImpl{psClientImpl, l, ""}
				if err := mockPsImpl.ExportToPubsub(ctx, csMap); err != nil {
					t.Fatalf("Testing Error, ErrorMessage: %v", err)
				}
			},
		},
		{
			name: "Pass to publish a message to pubsub with ordering key.",
			runner: func(t *testing.T) {
				psClientImpl := &mockPubsubClientImpl{pubsubClient: nil, cs: testCsArray}
				mockPsImpl := PubsubImpl{psClientImpl, l, "documentKey"}
				if err := mockPsImpl.ExportToPubsub(ctx, csMap); err != nil {
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

				psClientImpl := &mockPubsubClientImpl{pubsubClient: nil, cs: nil}
				mockPsImpl := PubsubImpl{psClientImpl, l, ""}
				if err := mockPsImpl.ExportToPubsub(ctx, csMap); err == nil {
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

				psClientImpl := &mockPubsubClientImpl{pubsubClient: nil, cs: nil}
				mockPsImpl := PubsubImpl{psClientImpl, l, ""}
				if err := mockPsImpl.ExportToPubsub(ctx, csMap); err == nil {
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

				psClientImpl := &mockPubsubClientImpl{pubsubClient: nil, cs: nil}
				mockPsImpl := PubsubImpl{psClientImpl, l, ""}
				if err := mockPsImpl.ExportToPubsub(ctx, csMap); err == nil {
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

				psClientImpl := &mockPubsubClientImpl{pubsubClient: nil, cs: nil}
				mockPsImpl := PubsubImpl{psClientImpl, l, ""}
				if err := mockPsImpl.ExportToPubsub(ctx, csMap); err == nil {
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

				psClientImpl := &mockPubsubClientImpl{pubsubClient: nil, cs: nil}
				mockPsImpl := PubsubImpl{psClientImpl, l, ""}
				if err := mockPsImpl.ExportToPubsub(ctx, csMap); err == nil {
					t.Fatalf("Not behaving as intended.")
				}
			},
		},
		{
			name: "Failed to get ordering key.",
			runner: func(t *testing.T) {
				psClientImpl := &mockPubsubClientImpl{pubsubClient: nil, cs: testCsArray}
				mockPsImpl := PubsubImpl{psClientImpl, l, "invalid-key"}
				if err := mockPsImpl.ExportToPubsub(ctx, csMap); err == nil {
					t.Fatalf("Not behaving as intended.")
				}
			},
		},
	}

	for _, v := range tests {
		t.Run(v.name, v.runner)
	}
}
