package main

import (
	"context"
	"go.uber.org/zap"
	"mxtransporter/application"
	"mxtransporter/config"
	"mxtransporter/pkg/client"
	"mxtransporter/pkg/logger"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var l *zap.SugaredLogger

	logCfg := config.LogConfig()
	l = logger.New(logCfg)

	mClient, err := client.NewMongoClient(ctx)
	if err != nil {
		l.Error(err)
		cancel()
	}
	defer mClient.Disconnect(ctx)

	watcherClient := &application.ChangeStremsWatcherClientImpl{mClient, application.ChangeStreamsExporterImpl{}}
	watcher := application.ChangeStremsWatcherImpl{watcherClient, l}

	if err := watcher.WatchChangeStreams(ctx); err != nil {
		l.Error(err)
		cancel()
	}
}
