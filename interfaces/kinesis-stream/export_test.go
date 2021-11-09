package kinesis_stream

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kinesis"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"reflect"
	"strconv"
	"testing"
)

// 重要
// kinesisifaceはaws sdk ver.2に対応していないので、使えない(ctxを関数の引数として挟めない)
// あとclientのインターフェースも用意されていないので、自分で定義する必要があり、その自分で定義したものをapplicationのfunctionでも利用する必要がある
// https://aws.github.io/aws-sdk-go-v2/docs/unit-testing/

// 目的の関数のモックを定義する
// mockPutRecordAPIがPutRecordのモック
// 実際にテストで検証する際に動作するのは、このモック
// モックの中身はcasesの中で入れてあげる
type mockPutRecordAPI func(ctx context.Context, params *kinesis.PutRecordInput, optFns ...func(*kinesis.Options)) (*kinesis.PutRecordOutput, error)

func (m mockPutRecordAPI) PutRecord(ctx context.Context, params *kinesis.PutRecordInput, optFns ...func(*kinesis.Options)) (*kinesis.PutRecordOutput, error) {
	return m(ctx, params, optFns...)
}

func Test_ExportToKinesisStream(t *testing.T) {
	csMap := primitive.M{
		"_id": primitive.M{"_data": "00000"},
		"operationType": "insert",
		"clusterTime": primitive.Timestamp{00000, 0},
		"fullDocument": primitive.M{"xxxxx": "xxxxx"},
		"ns": primitive.M{"xxxxx": "xxxxx"},
		"documentKey": primitive.M{"xxxxx": "xxxxx"},
		"updateDescription": primitive.M{"xxxxx": "xxxxx"},
	}

	j := "{\"_data\":\"00000\"}|insert|1970-01-01 09:00:00|{\"xxxxx\":\"xxxxx\"}|{\"xxxxx\":\"xxxxx\"}|{\"xxxxx\":\"xxxxx\"}|{\"xxxxx\":\"xxxxx\"}\n"
	rt := "00000"

	// casesで定義するのはExportToKinesisStreamのmock引数
	cases := []struct {
		client func(t *testing.T) kinesisPutRecordAPI
		cs primitive.M
	}{
		{
			//clientに関数を入れる
			// 返り値はkinesisClient
			client: func(t *testing.T) kinesisPutRecordAPI {
				//mockPutRecordAPIというclientのモックに対して中身(function)を入れてあげる
				// 実際にExportToKinesisStreamで呼ばれるPutRecord自体clientからよびだされているもので、
				// client内部の関数のモックをここで定義している
				// なので後でテストをrunするときにExportToKinesisStreamを呼び出すが
				// 内部のsdkはこのモックで定義した関数で動作する
				return mockPutRecordAPI(func(ctx context.Context, params *kinesis.PutRecordInput, optFns ...func(*kinesis.Options)) (*kinesis.PutRecordOutput, error) {
					t.Helper()
					//ここでparamsの値に対するアンチパターンを定義していく
					// TODO
					// ここのparamsがcsMapを参考にしていないと意味がなくないか？
					if params.Data == nil {
						t.Fatal("expect data to not be nil")
					}
					if e, a := []byte(j), params.Data; !reflect.DeepEqual(e, a) {
						t.Errorf("expect %v, got %v", e, a)
					}
					if params.PartitionKey == nil {
						t.Fatal("expect partition key to not be nil")
					}
					if e, a := aws.String(rt), params.PartitionKey; !reflect.DeepEqual(e, a) {
						t.Errorf("expect %v, got %v", e, a)
					}
					if params.StreamName == nil {
						t.Fatal("expect stream name to not be nil")
					}
					return nil, nil
				})
			},
			cs: csMap,
		},
	}

	for i, tt := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			ctx := context.TODO()
			err := ExportToKinesisStream(ctx, tt.cs, tt.client(t))
			if err != nil {
				t.Fatalf("expect no error, got %v", err)
			}
		})
	}
}

