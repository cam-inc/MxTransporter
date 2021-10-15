package main

import (
	"mxtransporter/application"
	mongoConnection "mxtransporter/interfaces/mongo"
	"context"
	"fmt"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client, err := mongoConnection.Connect(ctx)
	if err != nil {
		fmt.Println(err)
		cancel()
	}
	defer client.Disconnect(ctx)

	if err := application.WatchChangeStreams(ctx, client); err != nil {
		fmt.Println(err)
		cancel()
	}
}