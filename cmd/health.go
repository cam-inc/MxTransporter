package main

import (
	"context"
	"github.com/spf13/cobra"
	"go.mongodb.org/mongo-driver/mongo"
	mongoConnection "mxtransporter/interfaces/mongo"
	"mxtransporter/pkg/client"
	"mxtransporter/pkg/logger"
	"os"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	l := logger.New()

	mongoClient, err := client.NewMongoClient(ctx)
	if err != nil {
		l.Error(err)
		cancel()
	}
	defer mongoClient.Disconnect(ctx)

	c := &cobra.Command{}
	c.RunE = func(c *cobra.Command, args []string) error {
		return mongoDbConnectionCheck(ctx, mongoClient)
	}

	if err := c.Execute(); err != nil {
		os.Exit(2)
	}

	l.Info("Status OK.")
	os.Exit(0)
}

func mongoDbConnectionCheck(ctx context.Context, client *mongo.Client) error {
	if err := mongoConnection.Health(ctx, client); err != nil {
		return err
	}
	return nil
}
