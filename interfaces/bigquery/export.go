package bigquery

import (
	"cloud.google.com/go/bigquery"
	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson/primitive"
	config "mxtransporter/config/bigquery"
	"mxtransporter/pkg/errors"
	"time"
)

type changeStreamTableSchema struct {
	ID                string
	OperationType     string
	ClusterTime       time.Time
	FullDocument      string
	Ns                string
	DocumentKey       string
	UpdateDescription string
}

//type bigqueryClient interface {
//	Dataset(id string) datasetClient
//}
//
//type datasetClient interface {
//	Table(tableID string) tableClient
//}
//
//type tableClient interface {
//	Inserter() inserterClient
//}
//
//type inserterClient interface {
//	Put(ctx context.Context, src interface{}) (err error)
//}

type bigqueryClient interface {
	Dataset(id string) *bigquery.Dataset
}

type datasetClient interface {
	Table(tableID string) *bigquery.Table
}

type tableClient interface {
	Inserter() *bigquery.Inserter
}

type inserterClient interface {
	Put(ctx context.Context, src interface{}) (err error)
}

func ExportToBigquery(ctx context.Context, cs primitive.M, client bigqueryClient) error {
	bigqueryConfig := config.BigqueryConfig()

	id, _ := json.Marshal(cs["_id"])
	operationType := cs["operationType"].(string)
	clusterTime := cs["clusterTime"].(primitive.Timestamp).T
	fullDocument, _ := json.Marshal(cs["fullDocument"])
	ns, _ := json.Marshal(cs["ns"])
	documentKey, _ := json.Marshal(cs["documentKey"])
	updateDescription, _ := json.Marshal(cs["updateDescription"])

	test := client.Dataset(bigqueryConfig.DataSet)
	print(test)
	inserter := client.Dataset(bigqueryConfig.DataSet).Table(bigqueryConfig.Table).Inserter()
	csItems := []changeStreamTableSchema{
		{
			ID:                string(id),
			OperationType:     operationType,
			ClusterTime:       time.Unix(int64(clusterTime), 0),
			FullDocument:      string(fullDocument),
			Ns:                string(ns),
			DocumentKey:       string(documentKey),
			UpdateDescription: string(updateDescription),
		},
	}

	if err := inserter.Put(ctx, csItems); err != nil {
		return errors.InternalServerErrorBigqueryInsert.Wrap("Failed to insert record to Bigquery.", err)
	}

	return nil
}
