package opensearch

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"

	opensearchConfig "github.com/cam-inc/mxtransporter/config/opensearch"
	"github.com/cam-inc/mxtransporter/pkg/errors"
	"github.com/opensearch-project/opensearch-go/v3/opensearchapi"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var osCfg = opensearchConfig.OpenSearchConfig()

type (
	openSearchSingleIndexer interface {
		create(ctx context.Context, item *singleIndexerItem) error
		update(ctx context.Context, item *singleIndexerItem) error
		delete(ctx context.Context, item *singleIndexerItem) error
	}

	openSearchBulkIndexer interface {
		add(ctx context.Context, item *bulkIndexerItem) error
	}

	OpenSearchImpl struct {
		OpenSearchSingleIndexer openSearchSingleIndexer
		OpenSearchBulkIndexer   openSearchBulkIndexer
	}

	openSearchSingleIndexerImpl struct {
		osClient *opensearchapi.Client
	}

	singleIndexerItem struct {
		index      string
		documentID string
		body       primitive.M
	}

	openSearchBulkIndexerImpl struct {
		osClient *opensearchapi.Client
		idBufs   map[string][]*bytes.Buffer // The key is the object_id of the mongo doc.
		didFlush bool
		config   bulkIndexerConfig
	}

	bulkIndexerItem struct {
		indexActionType string
		index           string
		documentID      string
		body            primitive.M
	}

	bulkIndexerConfig struct {
		flushBytes         int
		aggregationEnabled bool
	}
)

// The return value, bool, indicates whether export was performed or not.
// If export was not performed due to buffering or an error, false is returned.
func (o *OpenSearchImpl) ExportToOpenSearch(ctx context.Context, cs primitive.M) (bool, error) {
	switch {
	case !osCfg.SyncEnabled && !osCfg.BulkEnabled:
		err := o.record(ctx, cs)
		if err != nil {
			return false, err
		}
		return true, nil
	case osCfg.SyncEnabled && !osCfg.BulkEnabled:
		err := o.sync(ctx, cs)
		if err != nil {
			return false, err
		}
		return true, nil
	case !osCfg.SyncEnabled && osCfg.BulkEnabled:
		err := o.bulkRecord(ctx, cs)
		didFlush := o.OpenSearchBulkIndexer.(*openSearchBulkIndexerImpl).didFlush
		if err != nil || !didFlush {
			return false, err
		}
		return true, nil
	case osCfg.SyncEnabled && osCfg.BulkEnabled:
		err := o.bulkSync(ctx, cs)
		didFlush := o.OpenSearchBulkIndexer.(*openSearchBulkIndexerImpl).didFlush
		if err != nil || !didFlush {
			return false, err
		}
		return true, nil
	}
	return false, errors.InternalServerError.New("opensearch config is invalid")
}

func (o *OpenSearchImpl) NewSingleIndexer(client *opensearchapi.Client) (*openSearchSingleIndexerImpl, error) {
	si := openSearchSingleIndexerImpl{
		osClient: client,
	}

	return &si, nil
}

func (o *OpenSearchImpl) record(ctx context.Context, cs primitive.M) error {
	siItem := &singleIndexerItem{
		index: osCfg.IndexName,
		body:  cs,
	}

	err := o.Create(ctx, siItem)
	if err != nil {
		return err
	}

	return nil
}

func (o *OpenSearchImpl) sync(ctx context.Context, cs primitive.M) error {
	opType := cs["operationType"].(string)
	index := osCfg.IndexName
	id := cs["documentKey"].(primitive.M)["_id"].(primitive.ObjectID).Hex()
	body := cs["fullDocument"].(primitive.M)

	siItem := &singleIndexerItem{
		index:      index,
		documentID: id,
		body:       body,
	}

	var err error
	switch opType {
	case "insert":
		err = o.Create(ctx, siItem)
	case "update":
		err = o.Update(ctx, siItem)
	case "delete":
		err = o.Delete(ctx, siItem)
	}
	if err != nil {
		return err
	}

	return nil
}

func (o *OpenSearchImpl) Create(ctx context.Context, siItem *singleIndexerItem) error {
	err := o.OpenSearchSingleIndexer.create(ctx, siItem)
	if err != nil {
		return err
	}
	return nil
}

func (si *openSearchSingleIndexerImpl) create(ctx context.Context, siItem *singleIndexerItem) error {
	jsonBody, err := json.Marshal(siItem.body)
	if err != nil {
		return errors.InternalServerErrorJsonMarshal.Wrap("Failed to json marshal body field of singleIndexerItem for create.", err)
	}

	indexReq := opensearchapi.IndexReq{
		Index: siItem.index,
		Body:  strings.NewReader(string(jsonBody)),
	}
	if osCfg.SyncEnabled {
		indexReq.DocumentID = siItem.documentID
	}

	_, err = si.osClient.Index(
		ctx,
		indexReq,
	)
	if err != nil {
		return errors.InternalServerErrorOpenSearchCreate.Wrap("Failed to create document to opensearch.", err)
	}

	return nil
}

func (o *OpenSearchImpl) Update(ctx context.Context, siItem *singleIndexerItem) error {
	err := o.OpenSearchSingleIndexer.update(ctx, siItem)
	if err != nil {
		return err
	}
	return nil
}

func (si *openSearchSingleIndexerImpl) update(ctx context.Context, siItem *singleIndexerItem) error {
	tempBody := map[string]interface{}{
		"doc": siItem.body,
	}

	jsonBody, err := json.Marshal(tempBody)
	if err != nil {
		return errors.InternalServerErrorJsonMarshal.Wrap("Failed to json marshal body field of singleIndexerItem for update.", err)
	}

	updateReq := opensearchapi.UpdateReq{
		Index:      siItem.index,
		DocumentID: siItem.documentID,
		Body:       strings.NewReader(string(jsonBody)),
	}

	_, err = si.osClient.Update(
		ctx,
		updateReq,
	)
	if err != nil {
		return errors.InternalServerErrorOpenSearchUpdate.Wrap("Failed to update document in opensearch.", err)
	}
	return nil
}

func (o *OpenSearchImpl) Delete(ctx context.Context, siItem *singleIndexerItem) error {
	err := o.OpenSearchSingleIndexer.delete(ctx, siItem)
	if err != nil {
		return err
	}
	return nil
}

func (si *openSearchSingleIndexerImpl) delete(ctx context.Context, siItem *singleIndexerItem) error {
	deleteReq := opensearchapi.DocumentDeleteReq{
		Index:      siItem.index,
		DocumentID: siItem.documentID,
	}

	_, err := si.osClient.Document.Delete(
		ctx,
		deleteReq,
	)
	if err != nil {
		return errors.InternalServerErrorOpenSearchDelete.Wrap("Failed to delete document in opensearch.", err)

	}
	return nil
}

func (o *OpenSearchImpl) NewBulkIndexer(ctx context.Context, client *opensearchapi.Client) (*openSearchBulkIndexerImpl, error) {
	biCfg := bulkIndexerConfig{
		flushBytes:         osCfg.BulkFlushBytes,
		aggregationEnabled: osCfg.BulkSyncAggregationEnabled,
	}

	bi := openSearchBulkIndexerImpl{
		osClient: client,
		idBufs:   make(map[string][]*bytes.Buffer),
		didFlush: false,
		config:   biCfg,
	}

	return &bi, nil
}

func (o *OpenSearchImpl) bulkRecord(ctx context.Context, cs primitive.M) error {
	indexActionType := "create"
	index := osCfg.IndexName
	id := cs["documentKey"].(primitive.M)["_id"].(primitive.ObjectID).Hex()

	biItem := &bulkIndexerItem{
		indexActionType: indexActionType,
		index:           index,
		documentID:      id, // Required as key for idBufs
		body:            cs,
	}

	return o.Add(ctx, biItem)
}

func (o *OpenSearchImpl) bulkSync(ctx context.Context, cs primitive.M) error {
	opType := cs["operationType"].(string)
	var indexActionType string
	// When buffering and then aggregating before indexing,
	// there is a possibility that the insert records of 'cs' might be skipped,
	// so it should be handled entirely with update (upsert).
	if opType == "insert" || opType == "update" {
		indexActionType = "update"
	} else if opType == "delete" {
		indexActionType = "delete"
	}
	index := osCfg.IndexName
	id := cs["documentKey"].(primitive.M)["_id"].(primitive.ObjectID).Hex()
	body := cs["fullDocument"].(primitive.M)

	biItem := &bulkIndexerItem{
		indexActionType: indexActionType,
		index:           index,
		documentID:      id,
		body:            body,
	}

	return o.Add(ctx, biItem)
}

func (item *bulkIndexerItem) writeMeta(buf *bytes.Buffer) error {
	meta := map[string]map[string]string{
		item.indexActionType: {
			"_index": item.index,
		},
	}
	if osCfg.SyncEnabled {
		meta[item.indexActionType]["DocumentID"] = item.documentID
	}

	jsonMeta, err := json.Marshal(meta)
	if err != nil {
		return errors.InternalServerErrorJsonMarshal.Wrap("Failed to marshal meta for opensearch.", err)
	}

	_, err = buf.Write(jsonMeta)
	if err != nil {
		return errors.InternalServerErrorOpenSearchAdd.Wrap("Failed to write meta to buffer.", err)
	}

	_, err = buf.WriteRune('\n')
	if err != nil {
		return errors.InternalServerErrorOpenSearchAdd.Wrap("Failed to write rune to buffer.", err)
	}

	return nil
}

func (item *bulkIndexerItem) writeBody(buf *bytes.Buffer) error {
	if item.body == nil {
		return nil
	}

	tempBody := make(map[string]interface{})
	if item.indexActionType == "update" {
		tempBody = map[string]interface{}{
			"doc": item.body,
		}
		if osCfg.SyncEnabled {
			tempBody["doc_as_upsert"] = true
		}
	} else if item.indexActionType == "create" {
		tempBody = item.body
	}

	jsonBody, err := json.Marshal(tempBody)
	if err != nil {
		return errors.InternalServerErrorJsonMarshal.Wrap("Failed to marshal body for opensearch.", err)
	}

	_, err = buf.Write(jsonBody)
	if err != nil {
		return errors.InternalServerErrorOpenSearchAdd.Wrap("Failed to write body to buffer.", err)
	}

	_, err = buf.WriteRune('\n')
	if err != nil {
		return errors.InternalServerErrorOpenSearchAdd.Wrap("Failed to write rune to buffer.", err)
	}

	return nil
}

func (bi *openSearchBulkIndexerImpl) flush(ctx context.Context) error {
	combinedBuf := new(bytes.Buffer)
	for _, bufs := range bi.idBufs {
		for _, buf := range bufs {
			_, err := combinedBuf.Write(buf.Bytes())
			if err != nil {
				return errors.InternalServerErrorOpenSearchAdd.Wrap("Failed to write bytes to combined buffer.", err)
			}
		}
	}

	req := opensearchapi.BulkReq{
		Index: osCfg.IndexName,
		Body:  combinedBuf,
	}

	_, err := bi.osClient.Bulk(ctx, req)
	if err != nil {
		return errors.InternalServerErrorOpenSearchAdd.Wrap("Failed to bulk index on opensearch.", err)
	}

	bi.idBufs = make(map[string][]*bytes.Buffer)

	return nil
}

func (o *OpenSearchImpl) Add(ctx context.Context, item *bulkIndexerItem) error {
	err := o.OpenSearchBulkIndexer.add(
		ctx,
		item,
	)
	if err != nil {
		return err
	}
	return nil
}

func (bi *openSearchBulkIndexerImpl) add(ctx context.Context, item *bulkIndexerItem) error {
	buf := new(bytes.Buffer)

	if err := item.writeMeta(buf); err != nil {
		return err
	}

	if err := item.writeBody(buf); err != nil {
		return err
	}

	if bi.config.aggregationEnabled {
		if _, ok := bi.idBufs[item.documentID]; ok {
			// Overwrite if present.
			bi.idBufs[item.documentID][len(bi.idBufs[item.documentID])-1] = buf
		} else {
			// If it does not exist, add it.
			bi.idBufs[item.documentID] = append(bi.idBufs[item.documentID], buf)
		}
	} else {
		// Add if not aggregated
		bi.idBufs[item.documentID] = append(bi.idBufs[item.documentID], buf)
	}

	// Flush if the total buffer size of the slices with the value of bi.idBufs exceeds flushBytes
	sumBufsSize := 0
	for _, bufs := range bi.idBufs {
		for _, buf := range bufs {
			sumBufsSize += buf.Len()
		}
	}
	if sumBufsSize >= bi.config.flushBytes {
		if err := bi.flush(ctx); err != nil {
			return err
		}
		bi.didFlush = true
	} else {
		bi.didFlush = false
	}

	return nil
}
