package storage

import (
	"bytes"
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/cam-inc/mxtransporter/pkg/client"
	"github.com/cam-inc/mxtransporter/pkg/errors"
	"io"
)

type (
	s3Cli struct {
		client *s3.Client
		bucket string
		region string
	}
)

func (s *s3Cli) GetObject(ctx context.Context, key string) ([]byte, error) {
	input := &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	}
	output, err := s.client.GetObject(ctx, input)
	if err != nil {
		return nil, errors.InternalServerErrorS3GetObject.Wrap("Failed to get object.", err)
	}
	defer output.Body.Close()
	return io.ReadAll(output.Body)
}

func (s *s3Cli) PutObject(ctx context.Context, key, value string) error {
	rtBuf := bytes.NewBuffer([]byte(value))
	input := &s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
		Body:   rtBuf,
	}
	_, err := s.client.PutObject(ctx, input)
	return errors.InternalServerErrorS3PutObject.Wrap("Failed to put object.", err)
}

func newS3(ctx context.Context, bucket, region string) (StorageClient, error) {
	cli := &s3Cli{}
	if s3Client, err := client.NewS3Client(ctx); err != nil {
		return nil, errors.InternalServerErrorS3NewClient.Wrap("Failed to initialize s3 client.", err)
	} else {
		cli.client = s3Client
	}
	cli.bucket = bucket
	cli.region = region
	return cli, nil
}
