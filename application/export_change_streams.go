package application

import (
	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/pubsub"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/kinesis"
	"github.com/cam-inc/mxtransporter/config"
	mongoConfig "github.com/cam-inc/mxtransporter/config/mongodb"
	interfaceForBigquery "github.com/cam-inc/mxtransporter/interfaces/bigquery"
	iff "github.com/cam-inc/mxtransporter/interfaces/file"
	interfaceForKinesisStream "github.com/cam-inc/mxtransporter/interfaces/kinesis-stream"
	mongoConnection "github.com/cam-inc/mxtransporter/interfaces/mongo"
	interfaceForPubsub "github.com/cam-inc/mxtransporter/interfaces/pubsub"
	"github.com/cam-inc/mxtransporter/pkg/client"
	"github.com/cam-inc/mxtransporter/pkg/common"
	"github.com/cam-inc/mxtransporter/pkg/errors"
	irt "github.com/cam-inc/mxtransporter/usecases/resume-token"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"strings"
	"time"
)

type agent string

const (
	BigQuery      agent = "bigquery"
	CloudPubSub   agent = "pubsub"
	KinesisStream agent = "kinesisStream"
	File          agent = "file"
)

type (
	changeStreamsWatcher interface {
		newBigqueryClient(ctx context.Context, projectID string) (*bigquery.Client, error)
		newPubsubClient(ctx context.Context, projectID string) (*pubsub.Client, error)
		newKinesisClient(ctx context.Context) (*kinesis.Client, error)
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

func (*ChangeStreamsWatcherClientImpl) newFileClient(_ context.Context) (iff.Exporter, error) {
	return iff.New(config.FileExportConfig()), nil
}

func (c *ChangeStreamsWatcherClientImpl) watch(ctx context.Context, ops *options.ChangeStreamOptions) (*mongo.ChangeStream, error) {
	mongoCfg := mongoConfig.MongoConfig()

	var pipeline mongo.Pipeline
	if mongoCfg.MongoWatchPipelineExcludeCsFulldocumentField != "" {
		excludeFields := strings.Split(mongoCfg.MongoWatchPipelineExcludeCsFulldocumentField, ",")

		if common.Contains(excludeFields, "_id") {
			return nil, errors.InternalServerError.New("The specified exclude field cannot be excluded from change streams.")
		} else if common.Contains(excludeFields, "operationType") {
			return nil, errors.InternalServerError.New("The specified exclude field cannot be excluded from change streams.")
		} else if common.Contains(excludeFields, "clusterTime") {
			return nil, errors.InternalServerError.New("The specified exclude field cannot be excluded from change streams.")
		} else if common.Contains(excludeFields, "fullDocument") {
			return nil, errors.InternalServerError.New("The specified exclude field cannot be excluded from change streams.")
		} else if common.Contains(excludeFields, "ns") {
			return nil, errors.InternalServerError.New("The specified exclude field cannot be excluded from change streams.")
		} else if common.Contains(excludeFields, "documentKey") {
			return nil, errors.InternalServerError.New("The specified exclude field cannot be excluded from change streams.")
		} else if common.Contains(excludeFields, "updateDescription") {
			return nil, errors.InternalServerError.New("The specified exclude field cannot be excluded from change streams.")
		}

		pipeline = mongo.Pipeline{bson.D{{"$unset", excludeFields}}}
	}

	cs, err := mongoConnection.Watch(ctx, c.MongoClient, pipeline, ops)
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
			bqClientImpl := &interfaceForBigquery.BigqueryClientImpl{bqClient}
			bqImpl = interfaceForBigquery.BigqueryImpl{bqClientImpl}
		case CloudPubSub:
			psClient, err := c.Watcher.newPubsubClient(ctx, projectID)
			if err != nil {
				return err
			}
			psClientImpl := &interfaceForPubsub.PubsubClientImpl{psClient, c.Log}
			psImpl = interfaceForPubsub.PubsubImpl{psClientImpl, c.Log}
		case KinesisStream:
			ksClient, err := c.Watcher.newKinesisClient(ctx)
			if err != nil {
				return err
			}
			ksClientImpl := &interfaceForKinesisStream.KinesisStreamClientImpl{ksClient}
			ksImpl = interfaceForKinesisStream.KinesisStreamImpl{ksClientImpl}
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
		exportToBigquery(ctx context.Context, cs primitive.M) error
		exportToPubsub(ctx context.Context, cs primitive.M) error
		exportToKinesisStream(ctx context.Context, cs primitive.M) error
		exportToFile(ctx context.Context, cs primitive.M) error
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

func (c *changeStreamsExporterClientImpl) exportToBigquery(ctx context.Context, cs primitive.M) error {
	return c.bq.ExportToBigquery(ctx, cs)
}

func (c *changeStreamsExporterClientImpl) exportToPubsub(ctx context.Context, cs primitive.M) error {
	return c.pubsub.ExportToPubsub(ctx, cs)
}

func (c *changeStreamsExporterClientImpl) exportToKinesisStream(ctx context.Context, cs primitive.M) error {
	return c.kinesisStream.ExportToKinesisStream(ctx, cs)
}

func (c *changeStreamsExporterClientImpl) exportToFile(ctx context.Context, cs primitive.M) error {
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

		var eg errgroup.Group
		for i := 0; i < len(expDstList); i++ {
			eDst := expDstList[i]
			eg.Go(func() error {
				switch agent(eDst) {
				case BigQuery:
					if err := c.exporter.exportToBigquery(ctx, csMap); err != nil {
						return err
					}
				case CloudPubSub:
					if err := c.exporter.exportToPubsub(ctx, csMap); err != nil {
						return err
					}
				case KinesisStream:
					if err := c.exporter.exportToKinesisStream(ctx, csMap); err != nil {
						return err
					}
				case File:
					if err := c.exporter.exportToFile(ctx, csMap); err != nil {
						return err
					}
				default:
					return errors.InternalServerError.Wrap("The export destination is wrong.", fmt.Errorf("you need to set the export destination in the environment variable correctly. you set %s", eDst))
				}
				return nil
			})
		}

		if err := eg.Wait(); err != nil {
			return err
		}

		csRt := csMap["_id"].(primitive.M)["_data"].(string)

		if err := c.exporter.saveResumeToken(ctx, csRt); err != nil {
			return err
		}
	}

	if err := c.exporter.err(); err != nil {
		return errors.InternalServerError.Wrap("Could not get the next event for change stream.", err)
	}

	c.log.Info("Acquisition of change streams was interrupted.")

	return nil
}
