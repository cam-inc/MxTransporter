package client

import (
	"context"
	"crypto/tls"
	"net/http"

	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/pubsub"
	"cloud.google.com/go/storage"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/kinesis"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	kinesisConfig "github.com/cam-inc/mxtransporter/config/kinesis-stream"
	mongoConfig "github.com/cam-inc/mxtransporter/config/mongodb"
	opensearchConfig "github.com/cam-inc/mxtransporter/config/opensearch"
	"github.com/cam-inc/mxtransporter/pkg/errors"
	opensearch "github.com/opensearch-project/opensearch-go/v3"
	"github.com/opensearch-project/opensearch-go/v3/opensearchapi"
	requestsigner "github.com/opensearch-project/opensearch-go/v3/signer/awsv2"
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

func NewOpenSearchClient(ctx context.Context) (*opensearchapi.Client, error) {
	osCfg := opensearchConfig.OpenSearchConfig()

	osClinetConfig := opensearch.Config{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Addresses: []string{osCfg.EndPoint},
	}

	if osCfg.UseAmazonOpenSearchService || osCfg.UseAmazonOpenSearchServerless {
		cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(osCfg.AwsRegion))
		if err != nil {
			return nil, err
		}

		var service string
		if osCfg.UseAmazonOpenSearchService {
			service = "es"
		} else {
			service = "aoss"
		}

		signer, err := requestsigner.NewSignerWithService(cfg, service)
		if err != nil {
			return nil, err
		}

		osClinetConfig.Signer = signer
	}

	c, err := opensearchapi.NewClient(opensearchapi.Config{Client: osClinetConfig})

	if err != nil {
		return nil, errors.InternalServerErrorClientGet.Wrap("failed aws load default config.", err)
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
