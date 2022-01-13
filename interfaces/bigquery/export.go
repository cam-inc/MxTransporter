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
	return b.BqClient.Dataset(dataset).Table(table).Inserter().Put(ctx, csItems)
}

func (b *BigqueryImpl) ExportToBigquery(ctx context.Context, cs primitive.M) error {
	bqCfg := bigqueryConfig.BigqueryConfig()

	id, err := json.Marshal(cs["_id"])
	if err != nil {
		return errors.InternalServerErrorJsonMarshal.Wrap("Failed to marshal change streams json.", err)
	}
	opType := cs["operationType"].(string)
	clusterTime := cs["clusterTime"].(primitive.Timestamp).T
	fullDoc, err := json.Marshal(cs["fullDocument"])
	if err != nil {
		return errors.InternalServerErrorJsonMarshal.Wrap("Failed to marshal change streams json.", err)
	}
	ns, err := json.Marshal(cs["ns"])
	if err != nil {
		return errors.InternalServerErrorJsonMarshal.Wrap("Failed to marshal change streams json.", err)
	}
	docKey, err := json.Marshal(cs["documentKey"])
	if err != nil {
		return errors.InternalServerErrorJsonMarshal.Wrap("Failed to marshal change streams json.", err)
	}
	updDesc, err := json.Marshal(cs["updateDescription"])
	if err != nil {
		return errors.InternalServerErrorJsonMarshal.Wrap("Failed to marshal change streams json.", err)
	}

	csItems := []ChangeStreamTableSchema{
		{
			ID:                string(id),
			OperationType:     opType,
			ClusterTime:       time.Unix(int64(clusterTime), 0),
			FullDocument:      string(fullDoc),
			Ns:                string(ns),
			DocumentKey:       string(docKey),
			UpdateDescription: string(updDesc),
		},
	}

	if err := b.Bq.putRecord(ctx, bqCfg.DataSet, bqCfg.Table, csItems); err != nil {
		return errors.InternalServerErrorBigqueryInsert.Wrap("Failed to insert record to Bigquery.", err)
	}

	return nil
}
