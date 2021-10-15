package main

import (
	mongoConnection "mxtransporter/interfaces/mongo"
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"go.mongodb.org/mongo-driver/mongo"
	"os"
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

	c := &cobra.Command{}
	c.RunE = func(c *cobra.Command, args []string) error {
		return mongoDbConnectionCheck(ctx, client)
	}

	if err := c.Execute(); err != nil {
		os.Exit(2)
	}

	fmt.Println("Status OK.")
	os.Exit(0)
}

func mongoDbConnectionCheck(ctx context.Context, client *mongo.Client) error {
	if err := mongoConnection.Health(ctx, client); err != nil {
		return err
	}
	return nil
}
