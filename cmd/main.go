package main

import (
	"context"
	"mxtransporter/application"
	"mxtransporter/pkg/client"
	"mxtransporter/pkg/logger"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	l := logger.New()

	mongoClient, err := client.NewMongoClient(ctx)
	if err != nil {
		l.ZLogger.Error(err)
		cancel()
	}
	defer mongoClient.Disconnect(ctx)

	watcherClient := &application.ChangeStremsWatcherClientImpl{mongoClient, application.ChangeStreamsExporterImpl{}}
	watcher := application.ChangeStremsWatcherImpl{watcherClient, l}

	if err := watcher.WatchChangeStreams(ctx); err != nil {
		l.ZLogger.Error(err)
		cancel()
	}
}
