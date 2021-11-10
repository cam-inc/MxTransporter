package bigquery

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"reflect"
	"strconv"
	"testing"
	"time"
)
// メモ
//メソッドチェーンされている箇所を外部functionにしてmockにしたら、テストの意味なくないか？？
// そこしか処理無いのに
//
// [[今回やればいいのは、csにmockデータを入れて、csItemsがそのmockから生成された値と同じかどうか]]
// そんで実際にfunctionのテストをするときにclientの関数もmockしないとtestが通らないよねって話
//
// 一旦保留

	//たけちゃのgoテストアドバイス

	//gomockを使うやり方
	//・gomockを使ってメソッドチェーンが表現できる
	//— mockを作るときに、そのmockの中でメソッドチェーンを表現して、強制的にチェーンされているfunctionが呼ばれるようにする
	//チェーンされている部分をラップするやり方

	//・ラップしたら、そのラップした部分のmockを作って、test対象のfunctionを呼ぶときに、引数として、そのfunctionを呼ぶようにする
	//(テスタブルな書き方)
	// DI

var csMap = primitive.M{
	"_id": primitive.M{"_data": "00000"},
	"operationType": "insert",
	"clusterTime": primitive.Timestamp{00000, 0},
	"fullDocument": primitive.M{"xxxxx": "xxxxx"},
	"ns": primitive.M{"xxxxx": "xxxxx"},
	"documentKey": primitive.M{"xxxxx": "xxxxx"},
	"updateDescription": primitive.M{"xxxxx": "xxxxx"},
}

type MockBigqueryFuncs struct {}

func (m *MockBigqueryFuncs) PutRecord(ctx context.Context, dataset string, table string, csItems []ChangeStreamTableSchema) error {
	_ = func (t *testing.T) error {
		t.Helper()
		testCsItems := []ChangeStreamTableSchema{
			{
				ID:                "{\"_data\":\"00000\"}",
				OperationType:     "insert",
				ClusterTime:       time.Unix(int64(csMap["clusterTime"].(primitive.Timestamp).T), 0),
				FullDocument:      "{\"xxxxx\":\"xxxxx\"}",
				Ns:                "{\"xxxxx\":\"xxxxx\"}",
				DocumentKey:       "{\"xxxxx\":\"xxxxx\"}",
				UpdateDescription: "{\"xxxxx\":\"xxxxx\"}",
			},
		}

		testDataset := "test dataset"
		testTable := "test table"


		dataset = testDataset
		table = testTable

		if dataset == "" {
			t.Fatal("expect dataset to not be nil")
		}
		if table == "" {
			t.Fatal("expect table to not be nil")
		}
		if csItems == nil {
			t.Fatal("expect csItems to not be nil")
		}
		if e, a := testCsItems, csItems; !reflect.DeepEqual(e, a) {
			t.Errorf("expect %v, got %v", e, a)
		}
		return nil
	}
	return nil
}


func Test_ExportToBigquery(t *testing.T) {
	cases := []struct {
		cs primitive.M
	}{
		{
			cs: csMap,
		},
	}

	function := &MockBigqueryFuncs{}

	for i, tt := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			ctx := context.TODO()
			err := ExportToBigquery(ctx, tt.cs, function)
			if err != nil {
				t.Fatalf("expect no error, got %v", err)
			}
		})
	}
}

