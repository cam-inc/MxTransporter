package kinesis_stream

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kinesis"
	"go.mongodb.org/mongo-driver/bson/primitive"
	kinesisConfig "mxtransporter/config/kinesis-stream"
	"mxtransporter/pkg/client"
	"mxtransporter/pkg/errors"
	"strings"
	"time"
)

var kinesisClient *kinesis.Client

type kinesisIf interface {
	PutRecord(ctx context.Context, streamName string, rt interface{}, csArray []string) error
}

type KinesisFuncs struct {}

func (k *KinesisFuncs) PutRecord(ctx context.Context, streamName string, rt interface{}, csArray []string) error {
	kinesisClient, err := client.NewKinesisClient(ctx)
	if err != nil {
		return err
	}

	_, err = kinesisClient.PutRecord(ctx, &kinesis.PutRecordInput{
		Data:         []byte(strings.Join(csArray, "|") + "\n"),
		PartitionKey: aws.String(rt.(string)),
		StreamName:   aws.String(streamName),
	})

	if err != nil {
		return errors.InternalServerErrorKinesisStreamPut.Wrap("Failed to put message into kinesis stream.", err)
	}

	return nil
}

func ExportToKinesisStream(ctx context.Context, cs primitive.M, ksif kinesisIf) error{
	kinesisStreamConfig := kinesisConfig.KinesisStreamConfig()

	rt := cs["_id"].(primitive.M)["_data"]

	id, _ := json.Marshal(cs["_id"])
	operationType := cs["operationType"].(string)
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


	if err := ksif.PutRecord(ctx, kinesisStreamConfig.StreamName, rt, r); err != nil {
		return errors.InternalServerErrorKinesisStreamPut.Wrap("Failed to put message into kinesis stream.", err)
	}

	return nil
}