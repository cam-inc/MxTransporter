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

	logConfig := config.LogConfig()
	l = logger.New(logConfig)

	mongoClient, err := client.NewMongoClient(ctx)
	if err != nil {
		l.Error(err)
		cancel()
	}
	defer mongoClient.Disconnect(ctx)

	watcherClient := &application.ChangeStremsWatcherClientImpl{mongoClient, application.ChangeStreamsExporterImpl{}}
	watcher := application.ChangeStremsWatcherImpl{watcherClient, l}

	if err := watcher.WatchChangeStreams(ctx); err != nil {
		l.Error(err)
		cancel()
	}
}
