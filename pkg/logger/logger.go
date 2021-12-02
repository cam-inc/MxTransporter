package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

type Log struct {
	Level           string
	Format          string
	OutputDirectory string
	OutputFile      string
}

func New(logCfg Log) *zap.SugaredLogger {
	var logFormat, logOutputPath, logOutputDir, logOutputF string

	level := zap.NewAtomicLevel()

	switch logCfg.Level {
	case "0":
		level.SetLevel(zapcore.InfoLevel)
	case "1":
		level.SetLevel(zapcore.ErrorLevel)
	}

	if logCfg.Format == "console" {
		logFormat = logCfg.Format
	} else {
		logFormat = "json"
	}

	if logCfg.OutputDirectory != "" && logCfg.OutputFile != "" {
		logOutputDir = logCfg.OutputDirectory
		logOutputF = logCfg.OutputFile
		if _, err := os.Stat(logOutputDir); os.IsNotExist(err) {
			os.MkdirAll(logOutputDir, 0777)
		}
		logOutputPath = logOutputDir + logOutputF
		os.OpenFile(logOutputPath, os.O_WRONLY|os.O_CREATE, 0664)
	} else {
		logOutputPath = "stdout"
	}

	config := zap.Config{
		Level:    level,
		Encoding: logFormat,
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:     "Time",
			LevelKey:    "Level",
			MessageKey:  "Msg",
			EncodeLevel: zapcore.CapitalLevelEncoder,
			EncodeTime:  zapcore.ISO8601TimeEncoder,
		},
		OutputPaths: []string{logOutputPath},
		// zap internal error output destination
		ErrorOutputPaths: []string{"stderr"},
	}
	logger, _ := config.Build()

	return logger.Sugar()
}
