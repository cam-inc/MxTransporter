package pubsub

import (
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

type mockPubsubFuncs struct{}

func (m *mockPubsubFuncs) PubsubTopic(_ context.Context, _ string) error {
	return nil
}

func (m *mockPubsubFuncs) PubsubSubscription(_ context.Context, _ string, _ string) error {
	return nil
}

func (m *mockPubsubFuncs) PublishMessage(_ context.Context, _ string, csArray []string) error {
	testCsArray := []string{
		`{"_data":"00000"}`,
		"insert",
		time.Unix(int64(csMap["clusterTime"].(primitive.Timestamp).T), 0).Format("2006-01-02 15:04:05"),
		`{"wwwww":"test full document"}`,
		`{"xxxxx":"test ns"}`,
		`{"yyyyy":"test document key"}`,
		`{"zzzzz":"test update description"}`,
	}

	if csArray == nil {
		return fmt.Errorf("expect csItems to not be nil")
	}
	if e, a := testCsArray, csArray; !reflect.DeepEqual(e, a) {
		return fmt.Errorf("expect %v, got %v", e, a)
	}
	return nil
}

func Test_ExportToPubSub(t *testing.T) {

	cases := []struct {
		cs       primitive.M
		function pubsubIf
	}{
		{
			cs:       csMap,
			function: &mockPubsubFuncs{},
		},
	}

	for i, tt := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			ctx := context.TODO()
			if err := ExportToPubSub(ctx, tt.cs, tt.function); err != nil {
				t.Fatalf("expect no error, got %v", err)
			}
		})
	}
}
