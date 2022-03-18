package main

import (
	"context"
	"github.com/cam-inc/mxtransporter/config"
	mongoConnection "github.com/cam-inc/mxtransporter/interfaces/mongo"
	"github.com/cam-inc/mxtransporter/pkg/client"
	"github.com/cam-inc/mxtransporter/pkg/logger"
	"github.com/spf13/cobra"
	"go.mongodb.org/mongo-driver/mongo"
	"os"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logCfg := config.LogConfig()
	l := logger.New(logCfg)

	mClient, err := client.NewMongoClient(ctx)
	if err != nil {
		l.Error(err)
		cancel()
	}
	defer mClient.Disconnect(ctx)

	c := &cobra.Command{}
	c.RunE = func(c *cobra.Command, args []string) error {
		return mongoDbConnectionCheck(ctx, mClient)
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
