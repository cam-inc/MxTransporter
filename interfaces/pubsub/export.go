package pubsub

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"cloud.google.com/go/pubsub"
	pubsubConfig "github.com/cam-inc/mxtransporter/config/pubsub"
	"github.com/cam-inc/mxtransporter/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

var results []*pubsub.PublishResult

type (
	IPubsub interface {
		topicExists(ctx context.Context, topicID string) (bool, error)
		createTopic(ctx context.Context, topicID string) (*pubsub.Topic, error)
		publishMessage(ctx context.Context, topicID string, csArray []string, pmo ...publishMessageOption) error
	}

	PubsubImpl struct {
		Pubsub     IPubsub
		Log        *zap.SugaredLogger
		OrderingBy string
	}

	PubsubClientImpl struct {
		PubsubClient *pubsub.Client
		Log          *zap.SugaredLogger
	}
)

func withOrderingKey(orderingKey string) publishMessageOption {
	return func(o *pubsub.Message) {
		o.OrderingKey = orderingKey
	}
}

type publishMessageOption func(opts *pubsub.Message)

func (p *PubsubClientImpl) topicExists(ctx context.Context, topicID string) (bool, error) {
	return p.PubsubClient.Topic(topicID).Exists(ctx)
}

func (p *PubsubClientImpl) createTopic(ctx context.Context, topicID string) (*pubsub.Topic, error) {
	return p.PubsubClient.CreateTopic(ctx, topicID)
}

func (p *PubsubClientImpl) publishMessage(ctx context.Context, topicID string, csArray []string, pmo ...publishMessageOption) error {
	topic := p.PubsubClient.Topic(topicID)
	defer topic.Stop()
	message := &pubsub.Message{
		Data: []byte(strings.Join(csArray, "|")),
	}
	for _, pmo := range pmo {
		pmo(message)
	}
	r := topic.Publish(ctx, message)

	for _, r := range append(results, r) {
		id, err := r.Get(ctx)
		if err != nil {
			return errors.InternalServerErrorPubSubPublish.Wrap("Failed to publish message.", err)
		}
		p.Log.Info("Published a message with a message ID: ", id)
	}

	return nil
}

// The return value, bool, indicates whether export was performed or not.
// If export was not performed due to buffering or an error, false is returned.
func (p *PubsubImpl) ExportToPubsub(ctx context.Context, cs primitive.M) (bool, error) {
	psCfg := pubsubConfig.PubSubConfig()

	topicID := psCfg.TopicName
	topicExistence, err := p.Pubsub.topicExists(ctx, topicID)
	if err != nil {
		return false, errors.InternalServerErrorPubSubFind.Wrap("Failed to check topic existence.", err)
	}
	if !topicExistence {
		p.Log.Info("Topic is not exists. Creating a topic.")

		var err error
		_, err = p.Pubsub.createTopic(ctx, topicID)
		if err != nil {
			return false, errors.InternalServerErrorPubSubCreate.Wrap("Failed to create topic.", err)
		}
		p.Log.Info("Successed to create topic. ")
	}

	id, err := json.Marshal(cs["_id"])
	if err != nil {
		return false, errors.InternalServerErrorJsonMarshal.Wrap("Failed to marshal change streams json _id parameter.", err)
	}
	opType := cs["operationType"].(string)
	clusterTime := cs["clusterTime"].(primitive.Timestamp).T
	fullDoc, err := json.Marshal(cs["fullDocument"])
	if err != nil {
		return false, errors.InternalServerErrorJsonMarshal.Wrap("Failed to marshal change streams json fullDocument parameter.", err)
	}
	ns, err := json.Marshal(cs["ns"])
	if err != nil {
		return false, errors.InternalServerErrorJsonMarshal.Wrap("Failed to marshal change streams json ns parameter.", err)
	}
	docKey, err := json.Marshal(cs["documentKey"])
	if err != nil {
		return false, errors.InternalServerErrorJsonMarshal.Wrap("Failed to marshal change streams json documentKey parameter.", err)
	}
	updDesc, err := json.Marshal(cs["updateDescription"])
	if err != nil {
		return false, errors.InternalServerErrorJsonMarshal.Wrap("Failed to marshal change streams json updateDescription parameter.", err)
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

	if p.OrderingBy != "" {
		key, err := p.orderingKey(cs)
		if err != nil {
			// TODO: error handling
			return false, err
		}
		err = p.Pubsub.publishMessage(ctx, topicID, r, withOrderingKey(key))
		if err != nil {
			// TODO: error handling
			return false, err
		}
		return true, nil
	}

	err = p.Pubsub.publishMessage(ctx, topicID, r)
	if err != nil {
		// TODO: error handling
		return false, err
	}
	return true, nil
}

func (p *PubsubImpl) orderingKey(cs primitive.M) (string, error) {
	key, ok := cs[p.OrderingBy]
	if !ok {
		return "", errors.InvalidErrorPubSubOrderingKey.New(fmt.Sprintf("Failed to get orderingKey cs: %v, orderingBy: %s", cs, p.OrderingBy))
	}
	return fmt.Sprintf("%v", key), nil
}
