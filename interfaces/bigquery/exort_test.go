package bigquery

import (
	"cloud.google.com/go/bigquery"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strconv"
	"testing"
)

//type mockBigqueryClient func(id string) datasetClient
//
//type mockDataset func(tableID string) tableClient
//
//type mockTable func() inserterClient
//
//type mockInserter func(ctx context.Context, src interface{}) (err error)

type mockBigqueryClient func(id string) *bigquery.Dataset

type mockDataset func(tableID string) *bigquery.Table

type mockTable func() *bigquery.Inserter

type mockInserter func(ctx context.Context, src interface{}) (err error)

//func (m mockBigqueryClient) Dataset(id string) datasetClient {
//	return m(id)
//}
//
//func (m mockDataset) Table(tableID string) tableClient {
//	return m(tableID)
//}
//
//func (m mockTable) Inserter() inserterClient {
//	return m()
//}
//
//func (m mockInserter) Put(ctx context.Context, src interface{}) (err error) {
//	return m(ctx, src)
//}

func (m mockBigqueryClient) Dataset(id string) *bigquery.Dataset {
	return m(id)
}

func (m mockDataset) Table(tableID string) *bigquery.Table {
	return m(tableID)
}

func (m mockTable) Inserter() *bigquery.Inserter {
	return m()
}

func (m mockInserter) Put(ctx context.Context, src interface{}) (err error) {
	return m(ctx, src)
}

func Test_ExportToBigquery(t *testing.T) {
	csMap := primitive.M{
		"_id": primitive.M{"_data": "00000"},
		"operationType": "insert",
		"clusterTime": primitive.Timestamp{00000, 0},
		"fullDocument": primitive.M{"xxxxx": "xxxxx"},
		"ns": primitive.M{"xxxxx": "xxxxx"},
		"documentKey": primitive.M{"xxxxx": "xxxxx"},
		"updateDescription": primitive.M{"xxxxx": "xxxxx"},
	}

	cases := []struct {
		client func(t *testing.T) bigqueryClient
		cs primitive.M
	}{
		{
			//client: func(t *testing.T) bigqueryClient {
			//	return mockBigqueryClient(func(id string) *bigquery.Dataset {
			//		t.Helper()
			//		id = "xxx"
			//		if id == "" {
			//			t.Fatal("expect id to not be nil")
			//		}
			//
			//		//sampleTbl := "xxx"
			//		//table := mockDataset(func(tableID string) *bigquery.Table {return nil}).Table(sampleTbl)
			//		//if table == nil {
			//		//	t.Fatal("expect table to not be nil")
			//		//}
			//		//
			//		//inserter := mockTable(func() *bigquery.Inserter {return nil}).Inserter()
			//		//if inserter == nil {
			//		//	t.Fatal("expect inserter to not be nil")
			//		//}
			//		//
			//		//ctx := context.TODO()
			//		//
			//		//src := ""
			//		//e := mockInserter(func(ctx context.Context, src interface{}) (err error) {return nil}).Put(ctx, src)
			//		//if e != nil {
			//		//	t.Fatal("error")
			//		//}
			//
			//		return nil
			//		//return mockDataset(func(tableID string) *bigquery.Table {
			//		//	tableID = "xxx"
			//		//	if tableID == "" {
			//		//		t.Fatal("expect tableID to not be nil")
			//		//	}
			//		//})
			//	})
			//},
			client: func(t *testing.T) bigqueryClient {
				return mockBigqueryClient(func(id string) *bigquery.Dataset {
					t.Helper()
					id = "xxx"
					if id == "" {
						t.Fatal("expect id to not be nil")
					}
					return mockDataset(func(tableID string) *bigquery.Table {
						tableID = "xxx"
						if tableID == "" {
							t.Fatal("expect tableID to not be nil")
						}
						return mockTable(func() *bigquery.Inserter {
							return mockInserter(func(ctx context.Context, src interface{}) (err error) {
								// TODO
								// srcの値を比較する時、test側で同じ値を作って比較する？？
								fmt.Printf("%v\n", src)
								if src == "" {
									t.Fatal("expect src to not be nil")
								}
								return nil
							})
						})
					})
				})
			},
			cs: csMap,
		},
	}

	for i, tt := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			ctx := context.TODO()
			err := ExportToBigquery(ctx, tt.cs, tt.client(t))
			if err != nil {
				t.Fatalf("expect no error, got %v", err)
			}
		})
	}
}

