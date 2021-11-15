package pubsub

import (
	"cloud.google.com/go/pubsub"
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"mxtransporter/config"
	pubsubConfig "mxtransporter/config/pubsub"
	"mxtransporter/pkg/client"
	"mxtransporter/pkg/errors"
	"strings"
	"time"
)

var pubSubClient *pubsub.Client
var gcpProjectID = config.FetchGcpProject().ProjectID

type pubsubIf interface {
	PubsubTopic(ctx context.Context, topicID string) error
	PubsubSubscription(ctx context.Context, topicID string, subscriptionID string) error
	PublishMessage(ctx context.Context, topicID string, csArray []string) error
}

type PubsubFuncs struct{}

func (p *PubsubFuncs) PubsubTopic(ctx context.Context, topicID string) error {
	pubSubClient, err := client.NewPubSubClient(ctx, gcpProjectID)
	if err != nil {
		return err
	}

	topic := pubSubClient.Topic(topicID)
	defer topic.Stop()

	topicExistence, err := topic.Exists(ctx)
	if err != nil {
		return errors.InternalServerErrorPubSubFind.Wrap("Failed to check topic existence.", err)
	}
	if topicExistence == false {
		fmt.Println("Topic is not exists. Creating a topic.")

		var err error
		_, err = pubSubClient.CreateTopic(ctx, topicID)
		if err != nil {
			return errors.InternalServerErrorPubSubCreate.Wrap("Failed to create topic.", err)
		}
		fmt.Println("Successed to create topic. ")
	}

	return nil
}

func (p *PubsubFuncs) PubsubSubscription(ctx context.Context, topicID string, subscriptionID string) error {
	pubSubClient, err := client.NewPubSubClient(ctx, gcpProjectID)
	if err != nil {
		return err
	}

	subscription := pubSubClient.Subscription(subscriptionID)

	subscriptionExistence, err := subscription.Exists(ctx)
	if err != nil {
		return errors.InternalServerErrorPubSubFind.Wrap("Failed to check subscription existence.", err)
	}
	if subscriptionExistence == false {
		fmt.Println("Subscription is not exists. Creating a subscription.")

		var err error
		_, err = pubSubClient.CreateSubscription(ctx, subscriptionID, pubsub.SubscriptionConfig{
			Topic:             pubSubClient.Topic(topicID),
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

func (p *PubsubFuncs) PublishMessage(ctx context.Context, topicID string, csArray []string) error {
	pubSubClient, err := client.NewPubSubClient(ctx, gcpProjectID)
	if err != nil {
		return err
	}

	topic := pubSubClient.Topic(topicID)
	defer topic.Stop()

	topic.Publish(ctx, &pubsub.Message{
		Data: []byte(strings.Join(csArray, "|")),
	})

	return nil
}

func ExportToPubSub(ctx context.Context, cs primitive.M, psif pubsubIf) error {
	pubSubConfig := pubsubConfig.PubSubConfig()

	topicID := pubSubConfig.MongoDbDatabase

	if err := psif.PubsubTopic(ctx, topicID); err != nil {
		return err
	}

	subscriptionID := pubSubConfig.MongoDbCollection

	if err := psif.PubsubSubscription(ctx, topicID, subscriptionID); err != nil {
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

	if err := psif.PublishMessage(ctx, topicID, r); err != nil {
		return err
	}

	return nil
}
