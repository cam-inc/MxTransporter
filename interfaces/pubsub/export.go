package pubsub

import (
	"cloud.google.com/go/pubsub"
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	config "mxtransporter/config/pubsub"
	"mxtransporter/pkg/errors"
	"strings"
	"time"
)

func ExportToPubSub(ctx context.Context, cs primitive.M, client *pubsub.Client) error {
	pubSubConfig := config.PubSubConfig()

	topicId := pubSubConfig.MongoDbDatabase

	var topic *pubsub.Topic
	topic = client.Topic(topicId)
	defer topic.Stop()

	topicExistence, err := topic.Exists(ctx)
	if err != nil {
		return errors.InternalServerErrorPubSub.Wrap("Failed to check topic existence.", err)
	}
	if topicExistence == false {
		fmt.Println("Topic is not exists. ")

		var err error
		topic, err = client.CreateTopic(ctx, topicId)
		if err != nil {
			return errors.InternalServerErrorPubSub.Wrap("Failed to create topic.", err)
		}
		fmt.Println("Successed to create topic. ")
	}

	subscriptionId := pubSubConfig.MongoDbCollection

	var subscription *pubsub.Subscription
	subscription = client.Subscription(subscriptionId)

	subscriptionExistence, err := subscription.Exists(ctx)
	if err != nil {
		return errors.InternalServerErrorPubSub.Wrap("Failed to check subscription existence.", err)
	}
	if subscriptionExistence == false {
		fmt.Println("Subscription is not exists. ")

		var err error
		subscription, err = client.CreateSubscription(ctx, subscriptionId, pubsub.SubscriptionConfig{
			Topic: topic,
			AckDeadline:       60 * time.Second,
			RetentionDuration: 24 * time.Hour,
		})
		if err != nil {
			return errors.InternalServerErrorPubSub.Wrap("Failed to create subscription.", err)
		}
		fmt.Println("Successed to create subscription. ")
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

	topic.Publish(ctx, &pubsub.Message{
		Data: []byte(strings.Join(r, "|")),
	})

	return nil
}