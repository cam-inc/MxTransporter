package main

import (
	"context"
	"github.com/cam-inc/mxtransporter/application"
	"github.com/cam-inc/mxtransporter/config"
	"github.com/cam-inc/mxtransporter/pkg/client"
	"github.com/cam-inc/mxtransporter/pkg/logger"
	"go.uber.org/zap"
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

	watcherClient := &application.ChangeStremsWatcherClientImpl{
		MongoClient: mClient,
		CsExporter:  application.ChangeStreamsExporterImpl{},
	}
	watcher := application.ChangeStremsWatcherImpl{
		Watcher: watcherClient,
		Log:     l,
	}

	if err := watcher.WatchChangeStreams(ctx); err != nil {
		l.Error(err)
		cancel()
	}
}
