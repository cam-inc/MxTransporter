package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	ZLogger *zap.SugaredLogger
}

var l Logger

func New() Logger {
	config := zap.Config{
		Level:       zap.NewAtomicLevel(),
		Development: false,
		Encoding:    "json",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:     "time",
			LevelKey:    "level",
			MessageKey:  "msg",
			EncodeLevel: zapcore.CapitalLevelEncoder,
			EncodeTime:  zapcore.ISO8601TimeEncoder,
		},
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}
	logger, _ := config.Build()

	l.ZLogger = logger.Sugar()
	return l
}
