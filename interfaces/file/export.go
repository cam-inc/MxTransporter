package file

import (
	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

type (
	Exporter interface {
		Export(ctx context.Context, cs primitive.M) error
	}

	ExporterConfig struct {
		LogType         string
		ChangeStreamKey string
		TimeKey         string
		NameKey         string
	}

	fileExporter struct {
		config *ExporterConfig
		log    *zap.Logger
	}
	timestamp struct {
		time.Time
		primitive.Timestamp
	}
	csDoc struct {
		ID                primitive.M `json:"_id"`
		OperationType     string      `json:"operationType"`
		ClusterTime       timestamp   `json:"clusterTime"`
		Ns                primitive.M `json:"ns"`
		FullDocument      primitive.M `json:"fullDocument"`
		DocumentKey       primitive.M `json:"documentKey"`
		UpdateDescription primitive.M `json:"updateDescription"`
	}
)

// MarshalJSON timestamp to time.Time (byte array)
func (t timestamp) MarshalJSON() ([]byte, error) {
	times := t.Time
	return json.Marshal(times)
}

// UnmarshalJSON primitive.Timestamp to timestamp struct
func (t *timestamp) UnmarshalJSON(data []byte) error {
	var pt primitive.Timestamp
	if err := json.Unmarshal(data, &pt); err != nil {
		return err
	}
	t.Timestamp = pt
	t.Time = time.Unix(int64(pt.T), 0)
	return nil
}

func (f *fileExporter) Export(_ context.Context, cs primitive.M) error {
	doc := &csDoc{}
	byteArray, err := json.Marshal(cs)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(byteArray, doc); err != nil {
		return err
	}

	f.log.Info("", zap.String("logType", f.config.LogType), zap.Any(f.config.ChangeStreamKey, doc))

	return nil
}

func New(cfg *ExporterConfig) Exporter {
	zconfig := zapcore.EncoderConfig{
		TimeKey:       cfg.TimeKey,
		NameKey:       cfg.NameKey,
		MessageKey:    zapcore.OmitKey,
		LevelKey:      zapcore.OmitKey,
		StacktraceKey: zapcore.OmitKey,
	}

	if cfg.LogType == "" {
		cfg.LogType = "changeStream"
	}
	if cfg.ChangeStreamKey == "" {
		cfg.ChangeStreamKey = "cs"
	}

	encoder := zapcore.NewJSONEncoder(zconfig)
	core := zapcore.NewCore(encoder, zapcore.WriteSyncer(os.Stdout), zap.NewAtomicLevelAt(zapcore.InfoLevel))
	log := zap.New(core)

	return &fileExporter{
		log:    log,
		config: cfg,
	}
}
