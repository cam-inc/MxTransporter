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
)

func (k *KinesisStreamClientImpl) putRecord(ctx context.Context, streamName string, rt interface{}, csArray []string) error {
	_, err := k.KinesisStreamClient.PutRecord(ctx, &kinesis.PutRecordInput{
		Data:         []byte(strings.Join(csArray, "|") + "\n"),
		PartitionKey: aws.String(rt.(string)),
		StreamName:   aws.String(streamName),
	})

	return err
}

func (k *KinesisStreamImpl) ExportToKinesisStream(ctx context.Context, cs primitive.M) error {
	ksCfg := kinesisConfig.KinesisStreamConfig()

	id, err := json.Marshal(cs["_id"])
	if err != nil {
		return errors.InternalServerErrorJsonMarshal.Wrap("Failed to marshal change streams json _id parameter.", err)
	}
	opType := cs["operationType"].(string)
	clusterTime := cs["clusterTime"].(primitive.Timestamp).T
	fullDoc, err := json.Marshal(cs["fullDocument"])
	if err != nil {
		return errors.InternalServerErrorJsonMarshal.Wrap("Failed to marshal change streams json fullDocument parameter.", err)
	}
	ns, err := json.Marshal(cs["ns"])
	if err != nil {
		return errors.InternalServerErrorJsonMarshal.Wrap("Failed to marshal change streams json ns parameter.", err)
	}
	docKey, err := json.Marshal(cs["documentKey"])
	if err != nil {
		return errors.InternalServerErrorJsonMarshal.Wrap("Failed to marshal change streams json documentKey parameter.", err)
	}
	updDesc, err := json.Marshal(cs["updateDescription"])
	if err != nil {
		return errors.InternalServerErrorJsonMarshal.Wrap("Failed to marshal change streams json updateDescription parameter.", err)
	}

	r := []string{
		string(id),
		opType,
		time.Unix(int64(clusterTime), 0).Format("2006-01-02 15:04:05"),
		string(fullDoc),
		string(ns),
		string(docKey),
		string(updDesc),
	}

	pm, ok := cs["_id"].(primitive.M)
	if !ok {
		return errors.InternalServerError.New("Failed to assert _id parameters of change streams.")
	}

	rt, exists := pm["_data"]
	if !exists {
		return errors.InternalServerError.New("Failed to get _data parameters of change streams.")
	}

	if err := k.KinesisStream.putRecord(ctx, ksCfg.StreamName, rt, r); err != nil {
		return errors.InternalServerErrorKinesisStreamPut.Wrap("Failed to put message into kinesis stream.", err)
	}

	return nil
}
