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

type (
	kinesisStreamClient interface {
		putRecord(ctx context.Context, streamName string, rt interface{}, csArray []string) error
	}

	KinesisStreamImpl struct {
		KinesisStream kinesisStreamClient
	}

	KinesisStreamClientImpl struct {
		KinesisStreamClient *kinesis.Client
	}

	mockKinesisStreamClientImpl struct {
		kinesisStreamClient *kinesis.Client
		rt                  string
		cs                  []string
	}
)

func (k *KinesisStreamClientImpl) putRecord(ctx context.Context, streamName string, rt interface{}, csArray []string) error {
	_, err := k.KinesisStreamClient.PutRecord(ctx, &kinesis.PutRecordInput{
		Data:         []byte(strings.Join(csArray, "|") + "\n"),
		PartitionKey: aws.String(rt.(string)),
		StreamName:   aws.String(streamName),
	})

	if err != nil {
		return errors.InternalServerErrorKinesisStreamPut.Wrap("Failed to put message into kinesis stream.", err)
	}

	return nil
}

func (k *KinesisStreamImpl) ExportToKinesisStream(ctx context.Context, cs primitive.M) error {
	kinesisStreamConfig := kinesisConfig.KinesisStreamConfig()

	rt := cs["_id"].(primitive.M)["_data"]

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

	r := []string{
		string(id),
		operationType,
		time.Unix(int64(clusterTime), 0).Format("2006-01-02 15:04:05"),
		string(fullDocument),
		string(ns),
		string(documentKey),
		string(updateDescription),
	}

	if err := k.KinesisStream.putRecord(ctx, kinesisStreamConfig.StreamName, rt, r); err != nil {
		return errors.InternalServerErrorKinesisStreamPut.Wrap("Failed to put message into kinesis stream.", err)
	}

	return nil
}
