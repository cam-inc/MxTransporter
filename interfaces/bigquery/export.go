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

type (
	bigqueryClient interface {
		putRecord(ctx context.Context, dataset string, table string, csItems []ChangeStreamTableSchema) error
	}

	BigqueryImpl struct {
		Bq bigqueryClient
	}

	BigqueryClientImpl struct {
		BqClient *bigquery.Client
	}
)

func (b *BigqueryClientImpl) putRecord(ctx context.Context, dataset string, table string, csItems []ChangeStreamTableSchema) error {
	if err := b.BqClient.Dataset(dataset).Table(table).Inserter().Put(ctx, csItems); err != nil {
		return errors.InternalServerErrorBigqueryInsert.Wrap("Failed to insert record to Bigquery.", err)
	}
	return nil
}

func (b *BigqueryImpl) ExportToBigquery(ctx context.Context, cs primitive.M) error {
	bqCfg := bigqueryConfig.BigqueryConfig()

	id, err := json.Marshal(cs["_id"])
	if err != nil {
		errors.InternalServerErrorJsonMarshal.Wrap("Failed to marshal json.", err)
	}
	operationType := cs["operationType"].(string)
	clusterTime := cs["clusterTime"].(primitive.Timestamp).T
	fullDocument, err := json.Marshal(cs["fullDocument"])
	if err != nil {
		errors.InternalServerErrorJsonMarshal.Wrap("Failed to marshal json.", err)
	}
	ns, err := json.Marshal(cs["ns"])
	if err != nil {
		errors.InternalServerErrorJsonMarshal.Wrap("Failed to marshal json.", err)
	}
	documentKey, err := json.Marshal(cs["documentKey"])
	if err != nil {
		errors.InternalServerErrorJsonMarshal.Wrap("Failed to marshal json.", err)
	}
	updateDescription, err := json.Marshal(cs["updateDescription"])
	if err != nil {
		errors.InternalServerErrorJsonMarshal.Wrap("Failed to marshal json.", err)
	}

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

	if err := b.Bq.putRecord(ctx, bqCfg.DataSet, bqCfg.Table, csItems); err != nil {
		return err
	}

	return nil
}
