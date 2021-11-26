package main

import (
	"context"
	"fmt"
	"mxtransporter/application"
	"mxtransporter/pkg/client"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mongoClient, err := client.NewMongoClient(ctx)
	if err != nil {
		fmt.Println(err)
		cancel()
	}
	defer mongoClient.Disconnect(ctx)

	// CsExporterの実態は後で入れます
	watcherClient := &application.ChangeStremsWatcherClientImpl{mongoClient, application.ChangeStreamsExporterImpl{}}
	watcher := application.ChangeStremsWatcherImpl{watcherClient}

	if err := watcher.WatchChangeStreams(ctx); err != nil {
		fmt.Println(err)
		cancel()
	}
}
