package pubsub

import (
	"cloud.google.com/go/pubsub"
	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson/primitive"
	pubsubConfig "mxtransporter/config/pubsub"
	"mxtransporter/pkg/errors"
	"mxtransporter/pkg/logger"
	"strings"
	"time"
)

type (
	pubsubClient interface {
		pubsubTopic(ctx context.Context, topicID string) error
		pubsubSubscription(ctx context.Context, topicID string, subscriptionID string) error
		publishMessage(ctx context.Context, topicID string, csArray []string) error
	}

	PubsubImpl struct {
		Pubsub pubsubClient
	}

	PubsubClientImpl struct {
		PubsubClient *pubsub.Client
		Log          logger.Logger
	}

	mockPubsubClientImpl struct {
		pubsubClient *pubsub.Client
		cs           []string
	}
)

func (p *PubsubClientImpl) pubsubTopic(ctx context.Context, topicID string) error {
	topic := p.PubsubClient.Topic(topicID)
	defer topic.Stop()

	topicExistence, err := topic.Exists(ctx)
	if err != nil {
		return errors.InternalServerErrorPubSubFind.Wrap("Failed to check topic existence.", err)
	}
	if topicExistence == false {
		p.Log.ZLogger.Info("Topic is not exists. Creating a topic.")

		var err error
		_, err = p.PubsubClient.CreateTopic(ctx, topicID)
		if err != nil {
			return errors.InternalServerErrorPubSubCreate.Wrap("Failed to create topic.", err)
		}
		p.Log.ZLogger.Info("Successed to create topic. ")
	}

	return nil
}

func (p *PubsubClientImpl) pubsubSubscription(ctx context.Context, topicID string, subscriptionID string) error {
	subscription := p.PubsubClient.Subscription(subscriptionID)

	subscriptionExistence, err := subscription.Exists(ctx)
	if err != nil {
		return errors.InternalServerErrorPubSubFind.Wrap("Failed to check subscription existence.", err)
	}
	if subscriptionExistence == false {
		p.Log.ZLogger.Info("Subscription is not exists. Creating a subscription.")

		var err error
		_, err = p.PubsubClient.CreateSubscription(ctx, subscriptionID, pubsub.SubscriptionConfig{
			Topic:             p.PubsubClient.Topic(topicID),
			AckDeadline:       60 * time.Second,
			RetentionDuration: 24 * time.Hour,
		})
		if err != nil {
			return errors.InternalServerErrorPubSubCreate.Wrap("Failed to create subscription.", err)
		}
		p.Log.ZLogger.Info("Successed to create subscription. ")
	}
	return nil
}

func (p *PubsubClientImpl) publishMessage(ctx context.Context, topicID string, csArray []string) error {
	topic := p.PubsubClient.Topic(topicID)
	defer topic.Stop()

	topic.Publish(ctx, &pubsub.Message{
		Data: []byte(strings.Join(csArray, "|")),
	})

	return nil
}

func (p *PubsubImpl) ExportToPubsub(ctx context.Context, cs primitive.M) error {
	pubSubConfig := pubsubConfig.PubSubConfig()

	topicID := pubSubConfig.MongoDbDatabase

	if err := p.Pubsub.pubsubTopic(ctx, topicID); err != nil {
		return err
	}

	subscriptionID := pubSubConfig.MongoDbCollection

	if err := p.Pubsub.pubsubSubscription(ctx, topicID, subscriptionID); err != nil {
		return err
	}

	id, err := json.Marshal(cs["_id"])
	if err != nil {
		errors.InternalServerErrorJsonMarshal.Wrap("Failed to marshal json.", err)
	}
	operationType := cs["operationType"].(string)
	clusterTime := cs["clusterTime"].(primitive.Timestamp).T
	fullDocument, err := json.Marshal(cs["fullDocument"])
	if err != nil {
		errors.InternalServerErrorJsonMarshal.Wrap("Failed to marshal json.", err)
	}
	ns, err := json.Marshal(cs["ns"])
	if err != nil {
		errors.InternalServerErrorJsonMarshal.Wrap("Failed to marshal json.", err)
	}
	documentKey, err := json.Marshal(cs["documentKey"])
	if err != nil {
		errors.InternalServerErrorJsonMarshal.Wrap("Failed to marshal json.", err)
	}
	updateDescription, err := json.Marshal(cs["updateDescription"])
	if err != nil {
		errors.InternalServerErrorJsonMarshal.Wrap("Failed to marshal json.", err)
	}

	r := []string{
		string(id),
		operationType,
		time.Unix(int64(clusterTime), 0).Format("2006-01-02 15:04:05"),
		string(fullDocument),
		string(ns),
		string(documentKey),
		string(updateDescription),
	}

	if err := p.Pubsub.publishMessage(ctx, topicID, r); err != nil {
		return err
	}

	return nil
}
