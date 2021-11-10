package pubsub

import (
	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson/primitive"
	config "mxtransporter/config/pubsub"
	"time"
)

type pubsubIf interface {
	PubsubTopic(ctx context.Context, topicID string) error
	PubsubSubscription(ctx context.Context, topicID string, subscriptionID string) error
    PublishMessage(ctx context.Context, topicID string, csArray []string) error
}

func ExportToPubSub(
		ctx context.Context,
		cs primitive.M,
		psif pubsubIf) error {
	pubSubConfig := config.PubSubConfig()

	topicID := pubSubConfig.MongoDbDatabase

	if err := psif.PubsubTopic(ctx, topicID); err != nil {
		return err
	}

	subscriptionID := pubSubConfig.MongoDbCollection

	if err := psif.PubsubSubscription(ctx, topicID, subscriptionID); err != nil {
		return err
	}

	id, _ := json.Marshal(cs["_id"])
	operationType := cs["operationType"].(string)
	clusterTime := cs["clusterTime"].(primitive.Timestamp).T
	fullDocument, _ := json.Marshal(cs["fullDocument"])
	ns, _ := json.Marshal(cs["ns"])
	documentKey, _ := json.Marshal(cs["documentKey"])
	updateDescription, _ := json.Marshal(cs["updateDescription"])

	r := []string{
		string(id),
		operationType,
		time.Unix(int64(clusterTime), 0).Format("2006-01-02 15:04:05"),
		string(fullDocument),
		string(ns),
		string(documentKey),
		string(updateDescription),
	}

	if err := psif.PublishMessage(ctx, topicID, r); err != nil {
		return err
	}

	return nil
}