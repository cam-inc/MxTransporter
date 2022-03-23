package storage

import (
	"cloud.google.com/go/storage"
	"context"
	"fmt"
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
	reader, _ := g.client.Bucket(g.bucket).Object(key).NewReader(ctx)
	defer reader.Close()
	return io.ReadAll(reader)
}

func (g *gcsCli) DeleteObject(ctx context.Context, key string) error {
	return g.client.Bucket(g.bucket).Object(key).Delete(ctx)
}

func (g *gcsCli) PutObject(ctx context.Context, key, value string) error {

	writer := g.client.Bucket(g.bucket).Object(key).NewWriter(ctx)
	defer writer.Close()
	_, err := writer.Write([]byte(value))
	return err
}

func NewGcs(ctx context.Context, bucket, region string) (StorageClient, error) {
	cli := &gcsCli{}
	gscCli, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	cli.client = gscCli
	cli.bucket = bucket
	cli.region = region
	fmt.Printf("DEBUG %v\n", cli)
	return cli, nil
}
