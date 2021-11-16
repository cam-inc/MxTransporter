package pubsub

import (
	"cloud.google.com/go/pubsub"
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	pubsubConfig "mxtransporter/config/pubsub"
	"mxtransporter/pkg/errors"
	"strings"
	"time"
)

type PubsubClient interface {
	PubsubTopic(ctx context.Context, topicID string, psClient *pubsub.Client) error
	PubsubSubscription(ctx context.Context, topicID string, subscriptionID string, psClient *pubsub.Client) error
	PublishMessage(ctx context.Context, topicID string, csArray []string, psClient *pubsub.Client) error
}

type PubsubClientImple struct {
	pubsubClient PubsubClient
}

func PubsubTopic(ctx context.Context, topicID string, psClient *pubsub.Client) error {
	topic := psClient.Topic(topicID)
	defer topic.Stop()

	topicExistence, err := topic.Exists(ctx)
	if err != nil {
		return errors.InternalServerErrorPubSubFind.Wrap("Failed to check topic existence.", err)
	}
	if topicExistence == false {
		fmt.Println("Topic is not exists. Creating a topic.")

		var err error
		_, err = psClient.CreateTopic(ctx, topicID)
		if err != nil {
			return errors.InternalServerErrorPubSubCreate.Wrap("Failed to create topic.", err)
		}
		fmt.Println("Successed to create topic. ")
	}

	return nil
}

func PubsubSubscription(ctx context.Context, topicID string, subscriptionID string, psClient *pubsub.Client) error {
	subscription := psClient.Subscription(subscriptionID)

	subscriptionExistence, err := subscription.Exists(ctx)
	if err != nil {
		return errors.InternalServerErrorPubSubFind.Wrap("Failed to check subscription existence.", err)
	}
	if subscriptionExistence == false {
		fmt.Println("Subscription is not exists. Creating a subscription.")

		var err error
		_, err = psClient.CreateSubscription(ctx, subscriptionID, pubsub.SubscriptionConfig{
			Topic:             psClient.Topic(topicID),
			AckDeadline:       60 * time.Second,
			RetentionDuration: 24 * time.Hour,
		})
		if err != nil {
			return errors.InternalServerErrorPubSubCreate.Wrap("Failed to create subscription.", err)
		}
		fmt.Println("Successed to create subscription. ")
	}
	return nil
}

func PublishMessage(ctx context.Context, topicID string, csArray []string, psClient *pubsub.Client) error {
	topic := psClient.Topic(topicID)
	defer topic.Stop()

	topic.Publish(ctx, &pubsub.Message{
		Data: []byte(strings.Join(csArray, "|")),
	})

	return nil
}

func NewPubsubClient(pubsubClient PubsubClient) *PubsubClientImple {
	return &PubsubClientImple{
		pubsubClient: pubsubClient,
	}
}

func (p *PubsubClientImple) ExportToPubSub(ctx context.Context, cs primitive.M, psClient *pubsub.Client) error {
	pubSubConfig := pubsubConfig.PubSubConfig()

	topicID := pubSubConfig.MongoDbDatabase

	if err := p.pubsubClient.PubsubTopic(ctx, topicID, psClient); err != nil {
		return err
	}

	subscriptionID := pubSubConfig.MongoDbCollection

	if err := p.pubsubClient.PubsubSubscription(ctx, topicID, subscriptionID, psClient); err != nil {
		return err
	}

	id, _ := json.Marshal(cs["_id"])
	operationType, _ := cs["operationType"].(string)
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

	if err := p.pubsubClient.PublishMessage(ctx, topicID, r, psClient); err != nil {
		return err
	}

	return nil
}
