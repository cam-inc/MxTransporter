package bigquery

import (
	"cloud.google.com/go/bigquery"
	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"mxtransporter/config"
	bigqueryConfig "mxtransporter/config/bigquery"
	"mxtransporter/pkg/client"
	"mxtransporter/pkg/errors"
	"time"
)

var bigqueryClient *bigquery.Client
var gcpProjectID = config.FetchGcpProject().ProjectID

type ChangeStreamTableSchema struct {
	ID                string
	OperationType     string
	ClusterTime       time.Time
	FullDocument      string
	Ns                string
	DocumentKey       string
	UpdateDescription string
}

type bigqueryIf interface {
	PutRecord(ctx context.Context, dataset string, table string, csItems []ChangeStreamTableSchema) error
}

type BigqueryFuncs struct {}

func (b *BigqueryFuncs) PutRecord(ctx context.Context, dataset string, table string, csItems []ChangeStreamTableSchema) error {
	bigqueryClient, err := client.NewBigqueryClient(ctx, gcpProjectID)
	if err != nil {
		return err
	}

	if err := bigqueryClient.Dataset(dataset).Table(table).Inserter().Put(ctx, csItems); err != nil {
		return errors.InternalServerErrorBigqueryInsert.Wrap("Failed to insert record to Bigquery.", err)
	}
	return nil
}

func ExportToBigquery(ctx context.Context, cs primitive.M, bqif bigqueryIf) error {
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

	if err := bqif.PutRecord(ctx, bigqueryConfig.DataSet, bigqueryConfig.Table, csItems); err != nil {
		return err
	}

	return nil
}

