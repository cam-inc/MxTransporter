//go:build test
// +build test

package pubsub

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
	"mxtransporter/config"
	"mxtransporter/pkg/logger"
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

func Test_ExportToPubSub(t *testing.T) {
	var l *zap.SugaredLogger
	logConfig := config.LogConfig()
	l = logger.New(logConfig)

	testCsArray := []string{
		`{"_data":"00000"}`,
		"insert",
		time.Unix(int64(csMap["clusterTime"].(primitive.Timestamp).T), 0).Format("2006-01-02 15:04:05"),
		`{"wwwww":"test full document"}`,
		`{"xxxxx":"test ns"}`,
		`{"yyyyy":"test document key"}`,
		`{"zzzzz":"test update description"}`,
	}

	t.Run("Test if the format of change streams works.", func(t *testing.T) {
		ctx := context.TODO()
		psClientImpl := &mockPubsubClientImpl{nil, testCsArray}
		mockPsImpl := PubsubImpl{psClientImpl, l}
		if err := mockPsImpl.ExportToPubsub(ctx, csMap); err != nil {
			t.Fatalf("Testing Error, ErrorMessage: %v", err)
		}
	})
}
