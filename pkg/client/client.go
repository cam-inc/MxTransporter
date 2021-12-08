package client

import (
	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/pubsub"
	"context"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/kinesis"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	kinesisConfig "mxtransporter/config/kinesis-stream"
	mongoConfig "mxtransporter/config/mongodb"
	"mxtransporter/pkg/errors"
)

func NewBigqueryClient(ctx context.Context, projectID string) (*bigquery.Client, error) {
	c, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		return nil, errors.InternalServerErrorClientGet.Wrap("bigquery client connection refused", err)
	}
	return c, nil
}

func NewPubsubClient(ctx context.Context, projectID string) (*pubsub.Client, error) {
	c, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return nil, errors.InternalServerErrorClientGet.Wrap("pubsub client connection refused", err)
	}
	return c, nil
}

func NewKinesisClient(ctx context.Context) (*kinesis.Client, error) {
	ksCfg := kinesisConfig.KinesisStreamConfig()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(ksCfg.KinesisStreamRegion))
	if err != nil {
		return nil, errors.InternalServerErrorClientGet.Wrap("aws client connection refused", err)
	}

	c := kinesis.NewFromConfig(cfg)

	return c, nil
}

func NewMongoClient(ctx context.Context) (*mongo.Client, error) {
	mongoCfg := mongoConfig.MongoConfig()
	c, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoCfg.MongoDbConnectionUrl))
	if err != nil {
		return nil, errors.InternalServerErrorMongoDbConnect.Wrap("mongodb connection refused.", err)
	}
	return c, nil
}
