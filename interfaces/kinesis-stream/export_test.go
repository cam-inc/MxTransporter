package kinesis_stream

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

type mockKinesisFuncs struct{}

func (m mockKinesisFuncs) PutRecord(_ context.Context, _ string, rt interface{}, csArray []string) error {
	testRt := "00000"

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
	if e, a := testRt, rt; !reflect.DeepEqual(e, a) {
		return fmt.Errorf("expect %v, got %v", e, a)
	}
	if e, a := testCsArray, csArray; !reflect.DeepEqual(e, a) {
		return fmt.Errorf("expect %v, got %v", e, a)
	}
	return nil
}

func Test_ExportToKinesisStream(t *testing.T) {
	cases := []struct {
		cs       primitive.M
		function kinesisIf
	}{
		{
			cs:       csMap,
			function: &mockKinesisFuncs{},
		},
	}

	for i, tt := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			ctx := context.TODO()
			if err := ExportToKinesisStream(ctx, tt.cs, tt.function); err != nil {
				t.Fatalf("expect no error, got %v", err)
			}
		})
	}
}
