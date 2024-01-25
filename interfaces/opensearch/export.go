package opensearch

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"
	"sync"
	"time"

	opensearchConfig "github.com/cam-inc/mxtransporter/config/opensearch"
	"github.com/cam-inc/mxtransporter/pkg/errors"
	"github.com/opensearch-project/opensearch-go/v3/opensearchapi"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var osCfg = opensearchConfig.OpenSearchConfig()

type (
	openSearchSingleIndexer interface {
		upsert(ctx context.Context, index string, id string, body primitive.M) error
		delete(ctx context.Context, index string, id string) error
	}

	openSearchBulkIndexer interface {
		add(ctx context.Context, item bulkIndexerItem) error
	}

	OpenSearchImpl struct {
		OpenSearchSingleIndexer openSearchSingleIndexer
		OpenSearchBulkIndexer   openSearchBulkIndexer
	}

	openSearchSingleIndexerImpl struct {
		osClient *opensearchapi.Client
	}

	openSearchBulkIndexerImpl struct {
		osClient *opensearchapi.Client
		queue    chan bulkIndexerItem
		idBufs   map[string][]*bytes.Buffer // キ-は mongo doc の object_id にする.
		ticker   *time.Ticker
		mu       sync.Mutex
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
		flushInterval      time.Duration
		aggregationEnabled bool
	}
)

func (o *OpenSearchImpl) ExportToOpenSearch(ctx context.Context, cs primitive.M) error {
	switch {
	case osCfg.SyncEnabled && osCfg.BulkEnabled:
		return o.bulkSync(ctx, cs)
	case osCfg.SyncEnabled && !osCfg.BulkEnabled:
		return o.sync(ctx, cs)
	case !osCfg.SyncEnabled && osCfg.BulkEnabled:
		return o.bulkRecord(ctx, cs)
	case !osCfg.SyncEnabled && !osCfg.BulkEnabled:
		return o.record(ctx, cs)
	default:
		return errors.InternalServerErrorMapKeyFind.New("invalid configuration: neither SyncEnabled nor BulkEnabled cases matched")
	}
}

func (o *OpenSearchImpl) record(ctx context.Context, cs primitive.M) error {
	opType := cs["operationType"].(string)
	index := osCfg.IndexName
	id := cs["documentKey"].(primitive.M)["_id"].(primitive.ObjectID).Hex()
	body := cs["fullDocument"].(primitive.M)

	switch opType {
	case "insert", "update":
		o.Upsert(ctx, index, id, body)
	case "delete":
		o.Delete(ctx, index, id)
	}
	return nil
}

func (o *OpenSearchImpl) sync(ctx context.Context, cs primitive.M) error {
	opType := cs["operationType"].(string)
	index := osCfg.IndexName
	id := cs["documentKey"].(primitive.M)["_id"].(primitive.ObjectID).Hex()
	body := cs["fullDocument"].(primitive.M)

	switch opType {
	case "insert", "update":
		o.Upsert(ctx, index, id, body)
	case "delete":
		o.Delete(ctx, index, id)
	}
	return nil
}

func (o *OpenSearchImpl) NewSingleIndexer(client *opensearchapi.Client) (*openSearchSingleIndexerImpl, error) {
	si := openSearchSingleIndexerImpl{
		osClient: client,
	}

	return &si, nil
}

func (o *OpenSearchImpl) Upsert(ctx context.Context, index string, id string, body primitive.M) error {
	err := o.OpenSearchSingleIndexer.upsert(ctx, index, id, body)
	if err != nil {
		return err

	}
	return nil
}

func (si *openSearchSingleIndexerImpl) upsert(ctx context.Context, index string, id string, body primitive.M) error {
	upsertBody := map[string]interface{}{
		"doc":           body,
		"doc_as_upsert": true,
	}

	jsonBody, err := json.Marshal(upsertBody)
	if err != nil {
		return err
	}
	_, err = si.osClient.Update(
		ctx,
		opensearchapi.UpdateReq{
			Index:      index,
			DocumentID: id,
			Body:       strings.NewReader(string(jsonBody)),
		},
	)
	if err != nil {
		return err
	}
	return nil
}

func (o *OpenSearchImpl) Delete(ctx context.Context, index string, id string) error {
	err := o.OpenSearchSingleIndexer.delete(ctx, index, id)
	if err != nil {
		return err

	}
	return nil
}

func (si *openSearchSingleIndexerImpl) delete(ctx context.Context, index string, id string) error {
	_, err := si.osClient.Document.Delete(
		ctx,
		opensearchapi.DocumentDeleteReq{
			Index:      index,
			DocumentID: id,
		},
	)
	if err != nil {
		return err

	}
	return nil
}

func (o *OpenSearchImpl) bulkRecord(ctx context.Context, cs primitive.M) error {
	var err error
	// 実装中

	return err
}

func (o *OpenSearchImpl) bulkSync(ctx context.Context, cs primitive.M) error {
	opType := cs["operationType"].(string)
	var indexActionType string
	// buffer を aggregate してから indexing する場合、insert が skip される可能性があるので全て upsert(update) で対応する。
	if opType == "insert" || opType == "update" {
		indexActionType = "update"
	} else if opType == "delete" {
		indexActionType = "delete"
	}
	index := osCfg.IndexName
	id := cs["documentKey"].(primitive.M)["_id"].(primitive.ObjectID).Hex()
	body := cs["fullDocument"].(primitive.M)

	biItem := bulkIndexerItem{
		indexActionType: indexActionType,
		index:           index,
		documentID:      id,
		body:            body,
	}

	return o.Add(ctx, biItem)
}

func (o *OpenSearchImpl) NewBulkIndexer(client *opensearchapi.Client) (*openSearchBulkIndexerImpl, error) {
	biCfg := bulkIndexerConfig{
		flushBytes:         osCfg.BulkFlushBytes,
		flushInterval:      time.Duration(osCfg.BulkFlushIntervalSeconds) * time.Second,
		aggregationEnabled: osCfg.BulkSyncAggregationEnabled,
	}

	bi := openSearchBulkIndexerImpl{
		osClient: client,
		config:   biCfg,
	}

	bi.init()

	return &bi, nil
}

func (bi *openSearchBulkIndexerImpl) init() {
	bi.queue = make(chan bulkIndexerItem)

	bi.run()

	bi.ticker = time.NewTicker(bi.config.flushInterval)

	go func() {
		ctx := context.Background()

		// range に修正
		for {
			select {
			case <-bi.ticker.C:
				bi.mu.Lock()
				if len(bi.idBufs) == 0 {
					bi.mu.Unlock()

					continue
				}
				if err := bi.flush(ctx); err != nil {
					bi.mu.Unlock()

					continue
				}
				bi.mu.Unlock()
			}
		}
	}()
}

func (bi *openSearchBulkIndexerImpl) run() {
	go func() {
		ctx := context.Background()
		for item := range bi.queue {
			bi.mu.Lock()

			buf := new(bytes.Buffer)

			if err := item.writeMeta(buf); err != nil {
				bi.mu.Unlock()
			}

			if err := item.writeBody(buf); err != nil {
				bi.mu.Unlock()
			}

			if bi.config.aggregationEnabled {
				if _, ok := bi.idBufs[item.documentID]; ok {
					// 存在したら上書きする
					bi.idBufs[item.documentID][len(bi.idBufs[item.documentID])-1] = buf
				} else {
					// 存在しなければ追加
					bi.idBufs[item.documentID] = append(bi.idBufs[item.documentID], buf)
				}
			} else {
				// 集約しない場合は追加
				bi.idBufs[item.documentID] = append(bi.idBufs[item.documentID], buf)
			}

			// bi.idBufs の 値のスライスが持つバッファーサイズの合計値が flushBytes を超えたら flush
			sumBufsSize := 0
			for _, bufs := range bi.idBufs {
				for _, buf := range bufs {
					sumBufsSize += buf.Len()
				}
			}
			if sumBufsSize >= bi.config.flushBytes {
				if err := bi.flush(ctx); err != nil {
					bi.mu.Unlock()
				}
			}
		}
	}()
}

func (item *bulkIndexerItem) writeMeta(buf *bytes.Buffer) error {
	meta := map[string]map[string]string{
		item.indexActionType: {
			"Index":      item.index,
			"DocumentID": item.documentID,
		},
	}

	jsonMeta, err := json.Marshal(meta)
	if err != nil {
		return err
	}

	_, err = buf.Write(jsonMeta)
	if err != nil {
		return err
	}

	_, err = buf.WriteRune('\n')
	if err != nil {
		return err
	}

	return nil
}

func (item *bulkIndexerItem) writeBody(buf *bytes.Buffer) error {
	if item.body == nil {
		return nil
	}

	upsertBody := map[string]interface{}{
		"doc":           item.body,
		"doc_as_upsert": true,
	}

	jsonBody, err := json.Marshal(upsertBody)
	if err != nil {
		return err
	}

	_, err = buf.Write(jsonBody)
	if err != nil {
		return err
	}

	_, err = buf.WriteRune('\n')
	if err != nil {
		return err
	}

	return nil
}

func (bi *openSearchBulkIndexerImpl) flush(ctx context.Context) error {
	combinedBuf := new(bytes.Buffer)
	for _, bufs := range bi.idBufs {
		for _, buf := range bufs {
			_, err := combinedBuf.Write(buf.Bytes())
			if err != nil {
				bi.mu.Unlock()
				return err
			}
		}
	}

	req := opensearchapi.BulkReq{
		Index: osCfg.IndexName,
		Body:  combinedBuf,
	}

	_, err := bi.osClient.Bulk(ctx, req)
	if err != nil {
		bi.mu.Unlock()
		return err
	}

	bi.idBufs = make(map[string][]*bytes.Buffer)

	bi.mu.Unlock()

	return nil
}

func (o *OpenSearchImpl) Add(ctx context.Context, item bulkIndexerItem) error {
	err := o.OpenSearchBulkIndexer.add(
		ctx,
		item,
	)
	if err != nil {
		return err
	}
	return nil
}

func (bi *openSearchBulkIndexerImpl) add(ctx context.Context, item bulkIndexerItem) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case bi.queue <- item:
	}

	return nil
}
