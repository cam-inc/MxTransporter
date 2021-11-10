package bigquery

import (
	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson/primitive"
	config "mxtransporter/config/bigquery"
	"time"
)

type bigqueryIf interface {
	PutRecord(ctx context.Context, dataset string, table string, csItems []ChangeStreamTableSchema) error
}

type ChangeStreamTableSchema struct {
	ID                string
	OperationType     string
	ClusterTime       time.Time
	FullDocument      string
	Ns                string
	DocumentKey       string
	UpdateDescription string
}

func ExportToBigquery(
		ctx context.Context,
		cs primitive.M,
		bqif bigqueryIf) error {
	bigqueryConfig := config.BigqueryConfig()

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

