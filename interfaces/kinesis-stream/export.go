package kinesis_stream

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kinesis"
	"go.mongodb.org/mongo-driver/bson/primitive"
	kinesisConfig "mxtransporter/config/kinesis-stream"
	"mxtransporter/pkg/errors"
	"strings"
	"time"
)

type KinesisClient interface {
	PutRecord(ctx context.Context, streamName string, rt interface{}, csArray []string, ksClient *kinesis.Client) error
}

type KinesisClientImpl struct {
	kinesisClient KinesisClient
}

func PutRecord(ctx context.Context, streamName string, rt interface{}, csArray []string, ksClient *kinesis.Client) error {
	_, err := ksClient.PutRecord(ctx, &kinesis.PutRecordInput{
		Data:         []byte(strings.Join(csArray, "|") + "\n"),
		PartitionKey: aws.String(rt.(string)),
		StreamName:   aws.String(streamName),
	})

	if err != nil {
		return errors.InternalServerErrorKinesisStreamPut.Wrap("Failed to put message into kinesis stream.", err)
	}

	return nil
}

func NewKinesisClient(kinesisClient KinesisClient) *KinesisClientImpl {
	return &KinesisClientImpl{
		kinesisClient: kinesisClient,
	}
}

func (k *KinesisClientImpl) ExportToKinesisStream(ctx context.Context, cs primitive.M, ksClient *kinesis.Client) error {
	kinesisStreamConfig := kinesisConfig.KinesisStreamConfig()

	rt := cs["_id"].(primitive.M)["_data"]

	id, _ := json.Marshal(cs["_id"])
	operationType, _ := cs["operationType"].(string)
	clusterTime := cs["clusterTime"].(primitive.Timestamp).T
	fullDocument, _ := json.Marshal(cs["fullDocument"])
	ns, _ := json.Marshal(cs["ns"])
	documentKey, _ := json.Marshal(cs["documentKey"])
	updateDescription, _ := json.Marshal(cs["updateDescription"])

	r := []string{
		string(id),
		operationType,
		time.Unix(int64(clusterTime), 0).Format("2006-01-02 15:04:05"),
		string(fullDocument),
		string(ns),
		string(documentKey),
		string(updateDescription),
	}

	if err := k.kinesisClient.PutRecord(ctx, kinesisStreamConfig.StreamName, rt, r, ksClient); err != nil {
		return err
	}

	return nil
}
