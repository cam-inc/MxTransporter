package storage

import (
	"cloud.google.com/go/storage"
	"context"
	"github.com/cam-inc/mxtransporter/pkg/client"
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
		return nil, err
	}
	defer reader.Close()
	o, err := io.ReadAll(reader)
	return o, err
}

func (g *gcsCli) PutObject(ctx context.Context, key, value string) error {
	writer := g.client.Bucket(g.bucket).Object(key).NewWriter(ctx)
	defer writer.Close()
	_, err := writer.Write([]byte(value))
	return err
}

func newGcs(ctx context.Context, bucket, region string) (StorageClient, error) {
	cli := &gcsCli{}
	gscCli, err := client.NewGcsClient(ctx)
	if err != nil {
		return nil, err
	}
	cli.client = gscCli
	cli.bucket = bucket
	cli.region = region
	return cli, nil
}
