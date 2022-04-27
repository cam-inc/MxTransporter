package storage

import (
	"cloud.google.com/go/storage"
	"context"
	"github.com/cam-inc/mxtransporter/pkg/client"
	"github.com/cam-inc/mxtransporter/pkg/errors"
	"io"
)

type (
	gcsCli struct {
		client *storage.Client
		bucket string
		region string
	}
)

func (g *gcsCli) GetObject(ctx context.Context, key string) ([]byte, error) {
	reader, err := g.client.Bucket(g.bucket).Object(key).NewReader(ctx)
	if err != nil {
		return nil, errors.InternalServerErrorGcsCreateNewReader.Wrap("Failed to create new reader.", err)
	}
	defer reader.Close()
	o, err := io.ReadAll(reader)
	return o, errors.InternalServerErrorGcsReader.Wrap("Failed to read object.", err)
}

func (g *gcsCli) PutObject(ctx context.Context, key, value string) error {
	writer := g.client.Bucket(g.bucket).Object(key).NewWriter(ctx)
	defer writer.Close()
	_, err := writer.Write([]byte(value))
	if err != nil {
		return errors.InternalServerErrorGcsWriteObject.Wrap("Failed to write object.", err)
	}
	return nil
}

func newGcs(ctx context.Context, bucket, region string) (StorageClient, error) {
	cli := &gcsCli{}
	gscCli, err := client.NewGcsClient(ctx)
	if err != nil {
		return nil, errors.InternalServerErrorGcsNewClient.Wrap("Failed to initialize gcs client.", err)
	}
	cli.client = gscCli
	cli.bucket = bucket
	cli.region = region
	return cli, nil
}
