package storage

import (
	"context"
	"fmt"
	"os"
)

type (
	fileStorageCli struct {
		volumePath string
	}
)

func (f *fileStorageCli) GetObject(ctx context.Context, key string) ([]byte, error) {

	rtByte, err := os.ReadFile(key)
	if err != nil {
		return nil, err
	}

	return rtByte, nil
}

func (f *fileStorageCli) DeleteObject(ctx context.Context, key string) error {
	return os.RemoveAll(key)
}

func (f *fileStorageCli) PutObject(ctx context.Context, key, value string) error {

	if _, err := os.Stat(f.volumePath); os.IsNotExist(err) {
		os.MkdirAll(f.volumePath, 0777)
	}

	fp, err := os.OpenFile(key, os.O_WRONLY|os.O_CREATE, 0664)

	if err != nil {
		return err
	}
	defer fp.Close()

	_, err = fp.WriteString(value)
	if err != nil {
		return err
	}
	return nil
}

func newLocalStorage(ctx context.Context, path string) (StorageClient, error) {
	cli := &fileStorageCli{
		volumePath: path,
	}
	fmt.Printf("DEBUG %v\n", cli)
	return cli, nil
}
