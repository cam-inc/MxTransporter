package file

import (
	"context"
	"encoding/json"
	"github.com/cam-inc/mxtransporter/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

type (
	Exporter interface {
		Put(ctx context.Context, cs primitive.M) error
	}
	fileExporter struct {
		log *zap.Logger
	}
)

func (f *fileExporter) Put(_ context.Context, cs primitive.M) error {
	data, err := unmarshalM(cs)
	if err != nil {
		return err
	}
	raw := json.RawMessage(data)
	f.log.Info("changeStream", zap.String("logType", "changeStream"), zap.Any("cs", &raw))
	return nil
}

func unmarshalM(cs primitive.M) ([]byte, error) {
	byteData, err := json.Marshal(cs)
	if err != nil {
		return nil, errors.InternalServerErrorJsonMarshal.Wrap("Failed Marshal primitive.M to JSON byte array", err)
	}
	return byteData, nil
}

func New(log *zap.Logger) Exporter {
	return &fileExporter{
		log: log.Named("exporter"),
	}
}
