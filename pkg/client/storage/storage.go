package storage

import (
	"context"
)

type (
	StorageClient interface {
		GetObject(ctx context.Context, key string) ([]byte, error)
		DeleteObject(ctx context.Context, key string) error
		PutObject(ctx context.Context, key, value string) error
	}

	serviceName string
)

const (
	s3Type    serviceName = "s3"
	gcsType   serviceName = "gcs"
	fileType  serviceName = "file"
	anonymous serviceName = "anonymous"
)

func ConvServiceName(name string) serviceName {
	switch serviceName(name) {
	case s3Type:
		return s3Type
	case gcsType:
		return gcsType
	case fileType:
		return fileType
	}
	return anonymous
}

func NewStorageClient(ctx context.Context, serviceName, path, bucketName, region string) (StorageClient, error) {
	switch ConvServiceName(serviceName) {
	case s3Type:
		return newS3(ctx, bucketName, region)
	case gcsType:
		return NewGcs(ctx, bucketName, region)
	case fileType:
		return newFile(ctx, path)
	}
	return newFile(ctx, path)
}
