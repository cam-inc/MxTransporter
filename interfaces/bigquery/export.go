package bigquery

import (
	"cloud.google.com/go/bigquery"
	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson/primitive"
	bigqueryConfig "mxtransporter/config/bigquery"
	"mxtransporter/pkg/errors"
	"time"
)

type ChangeStreamTableSchema struct {
	ID                string
	OperationType     string
	ClusterTime       time.Time
	FullDocument      string
	Ns                string
	DocumentKey       string
	UpdateDescription string
}

type BigqueryClient interface {
	PutRecord(ctx context.Context, dataset string, table string, csItems []ChangeStreamTableSchema, bqClient *bigquery.Client) error
}

type BigqueryClientImpl struct {
	bigqueryClient BigqueryClient
}

func PutRecord(ctx context.Context, dataset string, table string, csItems []ChangeStreamTableSchema, bqClient *bigquery.Client) error {
	if err := bqClient.Dataset(dataset).Table(table).Inserter().Put(ctx, csItems); err != nil {
		return errors.InternalServerErrorBigqueryInsert.Wrap("Failed to insert record to Bigquery.", err)
	}
	return nil
}

func NewBigqueryClient(bigqueryClient BigqueryClient) *BigqueryClientImpl {
	return &BigqueryClientImpl{
		bigqueryClient: bigqueryClient,
	}
}

func (b *BigqueryClientImpl) ExportToBigquery(ctx context.Context, cs primitive.M, bqClient *bigquery.Client) error {
	bigqueryConfig := bigqueryConfig.BigqueryConfig()

	id, _ := json.Marshal(cs["_id"])
	operationType := cs["operationType"].(string)
	clusterTime := cs["clusterTime"].(primitive.Timestamp).T
	fullDocument, _ := json.Marshal(cs["fullDocument"])
	ns, _ := json.Marshal(cs["ns"])
	documentKey, _ := json.Marshal(cs["documentKey"])
	updateDescription, _ := json.Marshal(cs["updateDescription"])

	csItems := []ChangeStreamTableSchema{
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

	if err := b.bigqueryClient.PutRecord(ctx, bigqueryConfig.DataSet, bigqueryConfig.Table, csItems, bqClient); err != nil {
		return err
	}

	return nil
}
