package elasticsearch

import (
	"context"
	"fmt"
	"strings"
	"time"

	elasticsearchConfig "github.com/cam-inc/mxtransporter/config/elasticsearch"
	"github.com/cam-inc/mxtransporter/pkg/errors"
	elasticsearch "github.com/elastic/go-elasticsearch/v8"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	elasticsearchClient interface {
		upsertFullDocument(ctx context.Context, index string, id string, fullDoc interface{}) error
		deleteFullDocument(ctx context.Context, index string, id string) error
		insertChangeStream(ctx context.Context, index string, cs ChangeStreamIndexSchema) error
	}

	ElasticsearchImpl struct {
		Es elasticsearchClient
	}

	ElasticsearchClientImpl struct {
		EsClient *elasticsearch.TypedClient
	}
)

func getString(m primitive.M, key string) (string, error) {
	val, ok := m[key].(string)
	if !ok {
		if m[key] == nil {
			return "", errors.InternalServerErrorMapKeyFind.New(fmt.Sprintf("Key does not exist or is nil, key: %s", key))
		}
		return "", errors.InternalServerErrorTypeAssertion.New(fmt.Sprintf("Failed to assert %s as string", key))
	}
	return val, nil
}

func getTimestamp(m primitive.M, key string) (primitive.Timestamp, error) {
	val, ok := m[key].(primitive.Timestamp)
	if !ok {
		if m[key] == nil {
			return primitive.Timestamp{}, errors.InternalServerErrorMapKeyFind.New(fmt.Sprintf("Key does not exist or is nil, key: %s", key))
		}
		return primitive.Timestamp{}, errors.InternalServerErrorTypeAssertion.New(fmt.Sprintf("Failed to assert %s as Timestamp", key))
	}
	return val, nil
}

func getMap(m primitive.M, key string) (primitive.M, error) {
	val, ok := m[key].(primitive.M)
	if !ok {
		if m[key] == nil {
			return nil, errors.InternalServerErrorMapKeyFind.New(fmt.Sprintf("Key does not exist or is nil, key: %s", key))
		}
		return nil, errors.InternalServerErrorTypeAssertion.New(fmt.Sprintf("Failed to assert %s as map", key))
	}
	return val, nil
}

func getObjectID(m primitive.M, key string) (primitive.ObjectID, error) {
	val, ok := m[key].(primitive.ObjectID)
	if !ok {
		if m[key] == nil {
			return primitive.ObjectID{}, errors.InternalServerErrorMapKeyFind.New(fmt.Sprintf("Key does not exist or is nil, key: %s", key))
		}
		return primitive.ObjectID{}, errors.InternalServerErrorTypeAssertion.New(fmt.Sprintf("Failed to assert %s as ObjectID", key))
	}
	return val, nil
}

func (e *ElasticsearchImpl) ExportToElasticsearch(ctx context.Context, cs primitive.M) error {
	esCfg := elasticsearchConfig.ElasticsearchConfig()

	if esCfg.SyncEnabled {
		return e.syncFullDocumentToElasticsearch(ctx, esCfg.IndexName, cs)
	}
	return e.insertChangeStreamToElasticsearch(ctx, esCfg.IndexName, cs)
}

func (e *ElasticsearchImpl) syncFullDocumentToElasticsearch(ctx context.Context, indexName string, cs primitive.M) error {
	opType, err := getString(cs, "operationType")
	if err != nil {
		return err
	}

	docKeyMap, err := getMap(cs, "documentKey")
	if err != nil {
		return err
	}

	objectID, err := getObjectID(docKeyMap, "_id")
	if err != nil {
		return err
	}

	id := objectID.Hex()

	if opType == "delete" {
		if err := e.Es.deleteFullDocument(ctx, indexName, id); err != nil {
			return errors.InternalServerErrorElasticsearchDelete.Wrap("Failed to delete document into Elasticsearch.", err)
		}
		return nil
	}

	fullDoc, err := getMap(cs, "fullDocument")
	if err != nil {
		return err
	}

	if err := e.Es.upsertFullDocument(ctx, indexName, id, fullDoc); err != nil {
		return errors.InternalServerErrorElasticsearchUpsert.Wrap("Failed to upsert document into Elasticsearch.", err)
	}

	return nil
}

func (e *ElasticsearchClientImpl) deleteFullDocument(ctx context.Context, index string, id string) error {
	_, err := e.EsClient.Delete(index, id).Do(ctx)
	return err
}

func (e *ElasticsearchClientImpl) upsertFullDocument(ctx context.Context, index string, id string, fullDoc interface{}) error {
	_, err := e.EsClient.Index(index).Request(fullDoc).Id(id).Do(ctx)
	return err
}

type ChangeStreamIndexSchema struct {
	ID                string
	OperationType     string
	ClusterTime       time.Time
	FullDocument      primitive.M
	Ns                primitive.M
	DocumentKey       primitive.M
	UpdateDescription primitive.M
}

func (e *ElasticsearchImpl) insertChangeStreamToElasticsearch(ctx context.Context, indexName string, cs primitive.M) error {
	csItem, err := createChangeStreamIndexSchema(cs)
	if err != nil {
		return errors.InternalServerErrorElasticsearchInsert.Wrap("Failed to create change stream index schema.", err)
	}
	if err := e.Es.insertChangeStream(ctx, indexName, csItem); err != nil {
		return errors.InternalServerErrorElasticsearchInsert.Wrap("Failed to insert change stream into Elasticsearch.", err)
	}

	return nil
}

func createChangeStreamIndexSchema(cs primitive.M) (ChangeStreamIndexSchema, error) {
	var schema ChangeStreamIndexSchema

	idData, err := getMap(cs, "_id")
	if err != nil {
		return schema, err
	}

	id, err := getString(idData, "_data")
	if err != nil {
		return schema, err
	}

	opType, err := getString(cs, "operationType")
	if err != nil {
		return schema, err
	}

	clusterTime, err := getTimestamp(cs, "clusterTime")
	if err != nil {
		return schema, err
	}

	ns, err := getMap(cs, "ns")
	if err != nil {
		return schema, err
	}

	docKey, err := getMap(cs, "documentKey")
	if err != nil {
		return schema, err
	}

	fullDocument, err := getMap(cs, "fullDocument")
	if err != nil {
		if strings.Contains(err.Error(), string(errors.InternalServerErrorMapKeyFind)) {
			return schema, nil
		}
		return schema, err
	}

	updateDescription, err := getMap(cs, "updateDescription")
	if err != nil {
		if strings.Contains(err.Error(), string(errors.InternalServerErrorMapKeyFind)) {
			return schema, nil
		}
		return schema, err
	}

	schema = ChangeStreamIndexSchema{
		ID:                id,
		OperationType:     opType,
		ClusterTime:       time.Unix(int64(clusterTime.T), 0),
		FullDocument:      fullDocument,
		Ns:                ns,
		DocumentKey:       docKey,
		UpdateDescription: updateDescription,
	}

	return schema, nil
}

func (e *ElasticsearchClientImpl) insertChangeStream(ctx context.Context, index string, cs ChangeStreamIndexSchema) error {
	_, err := e.EsClient.Index(index).Request(cs).Do(ctx)
	return err
}
