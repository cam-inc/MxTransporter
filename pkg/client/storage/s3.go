package storage

import (
	"bytes"
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"io"
)

type (
	s3client interface {
		PutObject(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error)
		GetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error)
		DeleteObject(ctx context.Context, params *s3.DeleteObjectInput, optFns ...func(*s3.Options)) (*s3.DeleteObjectOutput, error)
	}
	s3Cli struct {
		client s3client
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
		return nil, err
	}
	defer output.Body.Close()
	return io.ReadAll(output.Body)
}

func (s *s3Cli) DeleteObject(ctx context.Context, key string) error {
	input := &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	}
	_, err := s.client.DeleteObject(ctx, input)
	return err
}

func (s *s3Cli) PutObject(ctx context.Context, key, value string) error {
	rtBuf := bytes.NewBuffer([]byte(value))
	input := &s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
		Body:   rtBuf,
	}
	_, err := s.client.PutObject(ctx, input)
	return err
}

func newS3(ctx context.Context, bucket, region string) (StorageClient, error) {
	cli := &s3Cli{}
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return nil, err
	}
	cli.client = s3.NewFromConfig(cfg)
	cli.bucket = bucket
	cli.region = region
	return cli, nil
}
