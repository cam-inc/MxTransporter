package storage

import (
	"context"
	"github.com/cam-inc/mxtransporter/pkg/errors"
	"os"
)

type (
	fileStorageCli struct {
		volumePath string
	}
)

func (f *fileStorageCli) GetObject(_ context.Context, key string) ([]byte, error) {

	rtByte, err := os.ReadFile(key)
	if err != nil {
		return nil, errors.InternalServerError.Wrap("Failed to read file.", err)
	}

	return rtByte, nil
}

func (f *fileStorageCli) PutObject(_ context.Context, key, value string) error {

	if _, err := os.Stat(f.volumePath); os.IsNotExist(err) {
		os.MkdirAll(f.volumePath, 0777)
	}

	fp, err := os.OpenFile(key, os.O_WRONLY|os.O_CREATE, 0664)

	if err != nil {
		return errors.InternalServerError.Wrap("Failed to open file.", err)
	}
	defer fp.Close()

	_, err = fp.WriteString(value)
	if err != nil {
		return errors.InternalServerError.Wrap("Failed to write to file.", err)
	}
	return nil
}

func newFile(_ context.Context, path string) (StorageClient, error) {
	cli := &fileStorageCli{
		volumePath: path,
	}
	return cli, nil
}
