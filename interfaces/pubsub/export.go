package pubsub

import (
	"cloud.google.com/go/pubsub"
	"context"
	"encoding/json"
	pubsubConfig "github.com/cam-inc/mxtransporter/config/pubsub"
	"github.com/cam-inc/mxtransporter/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
	"strings"
	"time"
)

var results []*pubsub.PublishResult

type (
	IPubsub interface {
		topicExists(ctx context.Context, topicID string) (bool, error)
		createTopic(ctx context.Context, topicID string) (*pubsub.Topic, error)
		subscriptionExists(ctx context.Context, subscriptionID string) (bool, error)
		createSubscription(ctx context.Context, topicID string, subscriptionID string) (*pubsub.Subscription, error)
		publishMessage(ctx context.Context, topicID string, csArray []string) error
	}

	PubsubImpl struct {
		Pubsub IPubsub
		Log    *zap.SugaredLogger
	}

	PubsubClientImpl struct {
		PubsubClient *pubsub.Client
		Log          *zap.SugaredLogger
	}
)

func (p *PubsubClientImpl) topicExists(ctx context.Context, topicID string) (bool, error) {
	return p.PubsubClient.Topic(topicID).Exists(ctx)
}

func (p *PubsubClientImpl) createTopic(ctx context.Context, topicID string) (*pubsub.Topic, error) {
	return p.PubsubClient.CreateTopic(ctx, topicID)
}

func (p *PubsubClientImpl) subscriptionExists(ctx context.Context, subscriptionID string) (bool, error) {
	return p.PubsubClient.Subscription(subscriptionID).Exists(ctx)
}

func (p *PubsubClientImpl) createSubscription(ctx context.Context, topicID string, subscriptionID string) (*pubsub.Subscription, error) {
	return p.PubsubClient.CreateSubscription(ctx, subscriptionID, pubsub.SubscriptionConfig{
		Topic:             p.PubsubClient.Topic(topicID),
		AckDeadline:       60 * time.Second,
		RetentionDuration: 24 * time.Hour,
	})
}

func (p *PubsubClientImpl) publishMessage(ctx context.Context, topicID string, csArray []string) error {
	topic := p.PubsubClient.Topic(topicID)
	defer topic.Stop()

	r := topic.Publish(ctx, &pubsub.Message{
		Data: []byte(strings.Join(csArray, "|")),
	})

	for _, r := range append(results, r) {
		id, err := r.Get(ctx)
		if err != nil {
			return errors.InternalServerErrorPubSubPublish.Wrap("Failed to publish message.", err)
		}
		p.Log.Info("Published a message with a message ID: ", id)
	}

	return nil
}

func (p *PubsubImpl) ExportToPubsub(ctx context.Context, cs primitive.M) error {
	psCfg := pubsubConfig.PubSubConfig()

	topicID := psCfg.MongoDbDatabase
	topicExistence, err := p.Pubsub.topicExists(ctx, topicID)
	if err != nil {
		return errors.InternalServerErrorPubSubFind.Wrap("Failed to check topic existence.", err)
	}
	if !topicExistence {
		p.Log.Info("Topic is not exists. Creating a topic.")

		var err error
		_, err = p.Pubsub.createTopic(ctx, topicID)
		if err != nil {
			return errors.InternalServerErrorPubSubCreate.Wrap("Failed to create topic.", err)
		}
		p.Log.Info("Successed to create topic. ")
	}

	subscriptionID := psCfg.MongoDbCollection
	subscriptionExistence, err := p.Pubsub.subscriptionExists(ctx, subscriptionID)
	if err != nil {
		return errors.InternalServerErrorPubSubFind.Wrap("Failed to check subscription existence.", err)
	}

	if !subscriptionExistence {
		p.Log.Info("Subscription is not exists. Creating a subscription.")

		var err error
		_, err = p.Pubsub.createSubscription(ctx, topicID, subscriptionID)
		if err != nil {
			return errors.InternalServerErrorPubSubCreate.Wrap("Failed to create subscription.", err)
		}
		p.Log.Info("Successed to create subscription. ")
	}

	id, err := json.Marshal(cs["_id"])
	if err != nil {
		return errors.InternalServerErrorJsonMarshal.Wrap("Failed to marshal change streams json _id parameter.", err)
	}
	opType := cs["operationType"].(string)
	clusterTime := cs["clusterTime"].(primitive.Timestamp).T
	fullDoc, err := json.Marshal(cs["fullDocument"])
	if err != nil {
		return errors.InternalServerErrorJsonMarshal.Wrap("Failed to marshal change streams json fullDocument parameter.", err)
	}
	ns, err := json.Marshal(cs["ns"])
	if err != nil {
		return errors.InternalServerErrorJsonMarshal.Wrap("Failed to marshal change streams json ns parameter.", err)
	}
	docKey, err := json.Marshal(cs["documentKey"])
	if err != nil {
		return errors.InternalServerErrorJsonMarshal.Wrap("Failed to marshal change streams json documentKey parameter.", err)
	}
	updDesc, err := json.Marshal(cs["updateDescription"])
	if err != nil {
		return errors.InternalServerErrorJsonMarshal.Wrap("Failed to marshal change streams json updateDescription parameter.", err)
	}

	r := []string{
		string(id),
		opType,
		time.Unix(int64(clusterTime), 0).Format("2006-01-02 15:04:05"),
		string(fullDoc),
		string(ns),
		string(docKey),
		string(updDesc),
	}

	if err := p.Pubsub.publishMessage(ctx, topicID, r); err != nil {
		return err
	}

	return nil
}
