package kinesis_stream

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kinesis"
	"go.mongodb.org/mongo-driver/bson/primitive"
	config "mxtransporter/config/kinesis-stream"
	"mxtransporter/pkg/errors"
	"strings"
	"time"
)

type kinesisPutRecordAPI interface {
	PutRecord(ctx context.Context, params *kinesis.PutRecordInput, optFns ...func(*kinesis.Options)) (*kinesis.PutRecordOutput, error)
}

func ExportToKinesisStream(ctx context.Context, cs primitive.M, client kinesisPutRecordAPI) error{
	kinesisStreamConfig := config.KinesisStreamConfig()

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

	_, err := client.PutRecord(ctx, &kinesis.PutRecordInput{
		Data:         []byte(strings.Join(r, "|") + "\n"),
		PartitionKey: aws.String(rt.(string)),
		StreamName:   aws.String(kinesisStreamConfig.StreamName),
	})

	if err != nil {
		return errors.InternalServerErrorKinesisStreamPut.Wrap("Failed to put message into kinesis stream.", err)
	}

	return nil
}