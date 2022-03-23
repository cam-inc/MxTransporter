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

	storageCli struct {
		service serviceName
	}
)

const (
	s3Type    serviceName = "s3"
	gcsType   serviceName = "gcs"
	anonymous serviceName = "anonymous"
)

func (s *storageCli) GetObject(key string) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (s *storageCli) PutObject(key string) error {
	//TODO implement me
	panic("implement me")
}
func (s *storageCli) DeleteObject(key string) error {
	//TODO implement me
	panic("implement me")
}

func ConvServiceName(name string) serviceName {
	switch serviceName(name) {
	case s3Type:
		return s3Type
	case gcsType:
		return gcsType
	}
	return anonymous
}

func NewStorageClient(ctx context.Context, serviceName, path, bucketName, region string) (StorageClient, error) {

	switch ConvServiceName(serviceName) {
	case s3Type:
		return newS3(ctx, bucketName, region)
	case gcsType:
		return NewGcs(ctx, bucketName, region)
	}

	return newLocalStorage(ctx, path)
}
