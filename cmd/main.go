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

	if err := application.WatchChangeStreams(ctx, mongoClient); err != nil {
		fmt.Println(err)
		cancel()
	}
}
