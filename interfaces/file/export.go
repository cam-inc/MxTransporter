package file

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
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
)

func (f *fileExporter) Export(ctx context.Context, cs primitive.M) error {

	f.log.Info("", zap.String("logType", f.config.LogType), zap.Any(f.config.ChangeStreamKey, cs))

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
