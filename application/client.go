package application

import (
	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/pubsub"
	"context"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/kinesis"
	"mxtransporter/pkg/errors"
)

func NewBigqueryClient(ctx context.Context, projectID string) (*bigquery.Client, error) {
	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		return nil, errors.InternalServerErrorBigquery.Wrap("client connection refused", err)
	}
	return client, nil
}

func NewPubSubClient(ctx context.Context, projectID string) (*pubsub.Client, error) {
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return nil, errors.InternalServerErrorPubSub.Wrap("client connection refused", err)
	}
	return client, nil
}

//TODO
///////////////////////////////
// 実際にAWS GCPのコンテナ環境で権限を引き継げるのか確認
/////////////////////////////////
func NewKinesisClient(ctx context.Context, projectID string, region string) (*kinesis.Client, error) {
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithSharedConfigProfile(projectID),
		config.WithRegion(region),
	)
	if err != nil {
		return nil, errors.InternalServerErrorKinesisStream.Wrap("aws client connection refused", err)
	}

	client := kinesis.NewFromConfig(cfg)

	return client, nil
}