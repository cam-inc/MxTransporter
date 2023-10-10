package client

import (
	"context"

	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/pubsub"
	"cloud.google.com/go/storage"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/kinesis"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	elasticsearchConfig "github.com/cam-inc/mxtransporter/config/elasticsearch"
	kinesisConfig "github.com/cam-inc/mxtransporter/config/kinesis-stream"
	mongoConfig "github.com/cam-inc/mxtransporter/config/mongodb"
	"github.com/cam-inc/mxtransporter/pkg/errors"
	elasticsearch "github.com/elastic/go-elasticsearch/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
		return nil, errors.InternalServerErrorClientGet.Wrap("failed aws load default config.", err)
	}

	c := kinesis.NewFromConfig(cfg)

	return c, nil
}

func NewElasticsearchClient(ctx context.Context) (*elasticsearch.TypedClient, error) {
	esCfg := elasticsearchConfig.ElasticsearchConfig()
	cfg := elasticsearch.Config{
		Addresses: []string{
			esCfg.ElasticsearchConnectionUrl,
		},
	}
	c, err := elasticsearch.NewTypedClient(cfg)
	if err != nil {
		return nil, errors.InternalServerErrorClientGet.Wrap("elasticsearch client connection refused", err)
	}

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

func NewS3Client(ctx context.Context) (*s3.Client, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}
	return s3.NewFromConfig(cfg), nil
}
func NewGcsClient(ctx context.Context) (*storage.Client, error) {
	return storage.NewClient(ctx)
}
