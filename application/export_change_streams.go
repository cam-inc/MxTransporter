package application

import (
	"context"
	"fmt"
	"strings"
	"time"

	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/pubsub"
	"github.com/aws/aws-sdk-go-v2/service/kinesis"
	"github.com/cam-inc/mxtransporter/config"
	opensearchConfig "github.com/cam-inc/mxtransporter/config/opensearch"
	pconfig "github.com/cam-inc/mxtransporter/config/pubsub"
	interfaceForBigquery "github.com/cam-inc/mxtransporter/interfaces/bigquery"
	iff "github.com/cam-inc/mxtransporter/interfaces/file"
	interfaceForKinesisStream "github.com/cam-inc/mxtransporter/interfaces/kinesis-stream"
	mongoConnection "github.com/cam-inc/mxtransporter/interfaces/mongo"
	interfaceForOpenSearch "github.com/cam-inc/mxtransporter/interfaces/opensearch"
	interfaceForPubsub "github.com/cam-inc/mxtransporter/interfaces/pubsub"
	"github.com/cam-inc/mxtransporter/pkg/client"
	"github.com/cam-inc/mxtransporter/pkg/errors"
	irt "github.com/cam-inc/mxtransporter/usecases/resume-token"
	"github.com/opensearch-project/opensearch-go/v3/opensearchapi"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

type agent string

const (
	BigQuery      agent = "bigquery"
	CloudPubSub   agent = "pubsub"
	KinesisStream agent = "kinesisStream"
	OpenSearch    agent = "opensearch"
	File          agent = "file"
)

type (
	changeStreamsWatcher interface {
		newBigqueryClient(ctx context.Context, projectID string) (*bigquery.Client, error)
		newPubsubClient(ctx context.Context, projectID string) (*pubsub.Client, error)
		newKinesisClient(ctx context.Context) (*kinesis.Client, error)
		newOpenSearchClient(ctx context.Context) (*opensearchapi.Client, error)
		watch(ctx context.Context, ops *options.ChangeStreamOptions) (*mongo.ChangeStream, error)
		newFileClient(ctx context.Context) (iff.Exporter, error)
		setCsExporter(exporter ChangeStreamsExporterImpl)
		exportChangeStreams(ctx context.Context) error
	}

	ChangeStreamsWatcherImpl struct {
		Watcher            changeStreamsWatcher
		Log                *zap.SugaredLogger
		resumeTokenManager irt.ResumeToken
	}

	ChangeStreamsWatcherClientImpl struct {
		MongoClient *mongo.Client
		CsExporter  ChangeStreamsExporterImpl
	}
)

func (*ChangeStreamsWatcherClientImpl) newBigqueryClient(ctx context.Context, projectID string) (*bigquery.Client, error) {
	bqClient, err := client.NewBigqueryClient(ctx, projectID)
	if err != nil {
		return nil, err
	}
	return bqClient, nil
}

func (*ChangeStreamsWatcherClientImpl) newPubsubClient(ctx context.Context, projectID string) (*pubsub.Client, error) {
	psClient, err := client.NewPubsubClient(ctx, projectID)
	if err != nil {
		return nil, err
	}
	return psClient, nil
}

func (*ChangeStreamsWatcherClientImpl) newKinesisClient(ctx context.Context) (*kinesis.Client, error) {
	ksClient, err := client.NewKinesisClient(ctx)
	if err != nil {
		return nil, err
	}
	return ksClient, nil
}

func (*ChangeStreamsWatcherClientImpl) newOpenSearchClient(ctx context.Context) (*opensearchapi.Client, error) {
	osClient, err := client.NewOpenSearchClient(ctx)
	if err != nil {
		return nil, err
	}
	return osClient, nil
}

func (*ChangeStreamsWatcherClientImpl) newFileClient(_ context.Context) (iff.Exporter, error) {
	return iff.New(config.FileExportConfig()), nil
}

func (c *ChangeStreamsWatcherClientImpl) watch(ctx context.Context, ops *options.ChangeStreamOptions) (*mongo.ChangeStream, error) {
	cs, err := mongoConnection.Watch(ctx, c.MongoClient, ops)
	if err != nil {
		return nil, err
	}
	return cs, nil
}

func (c *ChangeStreamsWatcherClientImpl) setCsExporter(exporter ChangeStreamsExporterImpl) {
	c.CsExporter = exporter
}

func (c *ChangeStreamsWatcherImpl) setResumeTokenManager(resumeToken irt.ResumeToken) {
	c.resumeTokenManager = resumeToken
}
func (c *ChangeStreamsWatcherClientImpl) exportChangeStreams(ctx context.Context) error {
	return c.CsExporter.exportChangeStreams(ctx)
}

func (c *ChangeStreamsWatcherImpl) WatchChangeStreams(ctx context.Context) error {

	if c.resumeTokenManager == nil {
		rtImpl, err := irt.New(ctx, c.Log)
		if err != nil {
			return err
		}
		c.resumeTokenManager = rtImpl
	}

	rt := c.resumeTokenManager.ReadResumeToken(ctx)
	ops := options.ChangeStream().SetFullDocument(options.UpdateLookup)

	if len(rt) == 0 {
		c.Log.Info("File saved resume token in is not exists. Get from the current change streams.")
	} else {
		var rt interface{} = map[string]string{"_data": strings.TrimRight(rt, "\n")}

		ops.SetResumeAfter(rt)
	}

	cs, err := c.Watcher.watch(ctx, ops)
	if err != nil {
		return err
	}

	expDst, err := config.FetchExportDestination()
	if err != nil {
		return err
	}

	expDstList := strings.Split(expDst, ",")

	projectID, err := config.FetchGcpProject()
	if err != nil && (strings.Contains(expDst, string(BigQuery)) || strings.Contains(expDst, string(CloudPubSub))) {
		return err
	}

	var (
		bqImpl interfaceForBigquery.BigqueryImpl
		psImpl interfaceForPubsub.PubsubImpl
		ksImpl interfaceForKinesisStream.KinesisStreamImpl
		osImpl interfaceForOpenSearch.OpenSearchImpl
		fe     iff.Exporter
	)

	for i := 0; i < len(expDstList); i++ {
		eDst := expDstList[i]
		switch agent(eDst) {
		case BigQuery:
			bqClient, err := c.Watcher.newBigqueryClient(ctx, projectID)
			if err != nil {
				return err
			}
			bqClientImpl := &interfaceForBigquery.BigqueryClientImpl{BqClient: bqClient}
			bqImpl = interfaceForBigquery.BigqueryImpl{Bq: bqClientImpl}
		case CloudPubSub:
			psClient, err := c.Watcher.newPubsubClient(ctx, projectID)
			if err != nil {
				return err
			}
			psClientImpl := &interfaceForPubsub.PubsubClientImpl{PubsubClient: psClient, Log: c.Log}
			psImpl = interfaceForPubsub.PubsubImpl{Pubsub: psClientImpl, Log: c.Log, OrderingBy: pconfig.PubSubConfig().OrderingBy}
		case KinesisStream:
			ksClient, err := c.Watcher.newKinesisClient(ctx)
			if err != nil {
				return err
			}
			ksClientImpl := &interfaceForKinesisStream.KinesisStreamClientImpl{KinesisStreamClient: ksClient}
			ksImpl = interfaceForKinesisStream.KinesisStreamImpl{KinesisStream: ksClientImpl}
		case OpenSearch:
			osClient, err := c.Watcher.newOpenSearchClient(ctx)
			if err != nil {
				return err
			}

			osImpl = interfaceForOpenSearch.OpenSearchImpl{}

			osCfg := opensearchConfig.OpenSearchConfig()
			if osCfg.BulkEnabled {
				bi, err := osImpl.NewBulkIndexer(ctx, osClient)
				if err != nil {
					return err
				}
				osImpl.OpenSearchBulkIndexer = bi
			} else {
				si, err := osImpl.NewSingleIndexer(osClient)
				if err != nil {
					return err
				}
				osImpl.OpenSearchSingleIndexer = si
			}
		case File:
			fCli, err := c.Watcher.newFileClient(ctx)
			if err != nil {
				return err
			}
			fe = fCli
		default:
			return errors.InternalServerError.Wrap("The export destination is wrong.", fmt.Errorf("you need to set the export destination in the environment variable correctly. you set %s", eDst))
		}
	}

	exporterClient := &changeStreamsExporterClientImpl{
		cs:            cs,
		bq:            bqImpl,
		pubsub:        psImpl,
		kinesisStream: ksImpl,
		opensearch:    osImpl,
		fileExporter:  fe,
		resumeToken:   c.resumeTokenManager,
	}
	exporter := ChangeStreamsExporterImpl{
		exporter: exporterClient,
		log:      c.Log,
	}

	c.Watcher.setCsExporter(exporter)

	if err := c.Watcher.exportChangeStreams(ctx); err != nil {
		return err
	}

	return nil
}

type (
	changeStremsExporter interface {
		next(ctx context.Context) bool
		decode() (primitive.M, error)
		close(ctx context.Context) error
		exportToBigquery(ctx context.Context, cs primitive.M) (bool, error)
		exportToPubsub(ctx context.Context, cs primitive.M) (bool, error)
		exportToKinesisStream(ctx context.Context, cs primitive.M) (bool, error)
		exportToOpenSearch(ctx context.Context, cs primitive.M) (bool, error)
		exportToFile(ctx context.Context, cs primitive.M) (bool, error)
		saveResumeToken(ctx context.Context, rt string) error
		err() error
	}

	ChangeStreamsExporterImpl struct {
		exporter changeStremsExporter
		log      *zap.SugaredLogger
	}

	changeStreamsExporterClientImpl struct {
		cs            *mongo.ChangeStream
		bq            interfaceForBigquery.BigqueryImpl
		pubsub        interfaceForPubsub.PubsubImpl
		kinesisStream interfaceForKinesisStream.KinesisStreamImpl
		opensearch    interfaceForOpenSearch.OpenSearchImpl
		fileExporter  iff.Exporter
		resumeToken   irt.ResumeToken
	}
)

func (c *changeStreamsExporterClientImpl) next(ctx context.Context) bool {
	return c.cs.Next(ctx)
}

func (c *changeStreamsExporterClientImpl) decode() (primitive.M, error) {
	var csMap primitive.M

	if err := c.cs.Decode(&csMap); err != nil {
		return nil, errors.InternalServerError.Wrap("Failed to decode change stream.", err)
	}
	return csMap, nil
}

func (c *changeStreamsExporterClientImpl) close(ctx context.Context) error {
	return c.cs.Close(ctx)
}

func (c *changeStreamsExporterClientImpl) exportToBigquery(ctx context.Context, cs primitive.M) (bool, error) {
	return c.bq.ExportToBigquery(ctx, cs)
}

func (c *changeStreamsExporterClientImpl) exportToPubsub(ctx context.Context, cs primitive.M) (bool, error) {
	return c.pubsub.ExportToPubsub(ctx, cs)
}

func (c *changeStreamsExporterClientImpl) exportToKinesisStream(ctx context.Context, cs primitive.M) (bool, error) {
	return c.kinesisStream.ExportToKinesisStream(ctx, cs)
}

func (c *changeStreamsExporterClientImpl) exportToOpenSearch(ctx context.Context, cs primitive.M) (bool, error) {
	return c.opensearch.ExportToOpenSearch(ctx, cs)
}

func (c *changeStreamsExporterClientImpl) exportToFile(ctx context.Context, cs primitive.M) (bool, error) {
	return c.fileExporter.Export(ctx, cs)
}

func (c *changeStreamsExporterClientImpl) saveResumeToken(ctx context.Context, rt string) error {
	return c.resumeToken.SaveResumeToken(ctx, rt)
}

func (c *changeStreamsExporterClientImpl) err() error {
	return c.cs.Err()
}

func (c *ChangeStreamsExporterImpl) exportChangeStreams(ctx context.Context) error {
	defer c.exporter.close(ctx)

	expDst, err := config.FetchExportDestination()
	if err != nil {
		return err
	}
	expDstList := strings.Split(expDst, ",")

	for c.exporter.next(ctx) {

		csMap, err := c.exporter.decode()
		if err != nil {
			return err
		}

		csDb := csMap["ns"].(primitive.M)["db"].(string)
		csColl := csMap["ns"].(primitive.M)["coll"].(string)
		csOpType := csMap["operationType"].(string)
		csClusterTimeInt := time.Unix(int64(csMap["clusterTime"].(primitive.Timestamp).T), 0)

		c.log.Infof("Success to get change-streams, database: %s, collection: %s, operationType: %s, updateTime: %s", csDb, csColl, csOpType, csClusterTimeInt)

		saveRtFlags := make(chan bool, len(expDstList))
		var eg errgroup.Group
		for i := 0; i < len(expDstList); i++ {
			eDst := expDstList[i]
			eg.Go(func() error {
				switch agent(eDst) {
				case BigQuery:
					if saveRtFlag, err := c.exporter.exportToBigquery(ctx, csMap); err != nil {
						return err
					} else {
						saveRtFlags <- saveRtFlag
					}
				case CloudPubSub:
					if saveRtFlag, err := c.exporter.exportToPubsub(ctx, csMap); err != nil {
						return err
					} else {
						saveRtFlags <- saveRtFlag
					}
				case KinesisStream:
					if saveRtFlag, err := c.exporter.exportToKinesisStream(ctx, csMap); err != nil {
						return err
					} else {
						saveRtFlags <- saveRtFlag
					}
				case OpenSearch:
					if saveRtFlag, err := c.exporter.exportToOpenSearch(ctx, csMap); err != nil {
						return err
					} else {
						saveRtFlags <- saveRtFlag
					}
				case File:
					if saveRtFlag, err := c.exporter.exportToFile(ctx, csMap); err != nil {
						return err
					} else {
						saveRtFlags <- saveRtFlag
					}
				default:
					return errors.InternalServerError.Wrap("The export destination is wrong.", fmt.Errorf("you need to set the export destination in the environment variable correctly. you set %s", eDst))
				}
				return nil
			})
		}

		if err := eg.Wait(); err != nil {
			close(saveRtFlags)
			return err
		}
		close(saveRtFlags)

		// Save only when all flags are true.
		skipSave := false
		for saveRtFlag := range saveRtFlags {
			if !saveRtFlag {
				skipSave = true
				break
			}
		}

		if !skipSave {
			csRt := csMap["_id"].(primitive.M)["_data"].(string)
			if err := c.exporter.saveResumeToken(ctx, csRt); err != nil {
				return err
			}
		}
	}

	if err := c.exporter.err(); err != nil {
		return errors.InternalServerError.Wrap("Could not get the next event for change stream.", err)
	}

	c.log.Info("Acquisition of change streams was interrupted.")

	return nil
}
