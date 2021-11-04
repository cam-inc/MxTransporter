package application

import (
	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/pubsub"
	"context"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/kinesis"
	kinesisConfig "mxtransporter/config/kinesis-stream"
	"mxtransporter/pkg/errors"
)

func NewBigqueryClient(ctx context.Context, projectID string) (*bigquery.Client, error) {
	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		return nil, errors.InternalServerErrorClientGet.Wrap("bigquery client connection refused", err)
	}
	return client, nil
}

func NewPubSubClient(ctx context.Context, projectID string) (*pubsub.Client, error) {
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return nil, errors.InternalServerErrorClientGet.Wrap("pubsub client connection refused", err)
	}
	return client, nil
}

func NewKinesisClient(ctx context.Context) (*kinesis.Client, error) {
	kinesisStreamConfig := kinesisConfig.KinesisStreamConfig()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(kinesisStreamConfig.KinesisStreamRegion))
	if err != nil {
		return nil, errors.InternalServerErrorClientGet.Wrap("aws client connection refused", err)
	}

	client := kinesis.NewFromConfig(cfg)

	return client, nil
}