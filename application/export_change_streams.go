package application

import (
	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/pubsub"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/kinesis"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"mxtransporter/config"
	interfaceForBigquery "mxtransporter/interfaces/bigquery"
	interfaceForKinesisStream "mxtransporter/interfaces/kinesis-stream"
	mongoConnection "mxtransporter/interfaces/mongo"
	interfaceForPubsub "mxtransporter/interfaces/pubsub"
	"mxtransporter/pkg/client"
	"mxtransporter/pkg/common"
	"mxtransporter/pkg/errors"
	interfaceForResumeToken "mxtransporter/usecases/resume-token"
	"os"
	"strings"
	"time"
)

type agent string

const (
	BigQuery      agent = "bigquery"
	CloudPubSub   agent = "pubsub"
	KinesisStream agent = "kinesisStream"
)

type (
	changeStremsWatcher interface {
		newBigqueryClient(ctx context.Context, projectID string) (*bigquery.Client, error)
		newPubsubClient(ctx context.Context, projectID string) (*pubsub.Client, error)
		newKinesisClient(ctx context.Context) (*kinesis.Client, error)
		watch(ctx context.Context, ops *options.ChangeStreamOptions) (*mongo.ChangeStream, error)
		setCsExporter(exporter ChangeStreamsExporterImpl)
		exportChangeStreams(ctx context.Context) error
	}

	ChangeStremsWatcherImpl struct {
		Watcher changeStremsWatcher
		Log     *zap.SugaredLogger
	}

	ChangeStremsWatcherClientImpl struct {
		MongoClient *mongo.Client
		CsExporter  ChangeStreamsExporterImpl
	}
)

func (_ *ChangeStremsWatcherClientImpl) newBigqueryClient(ctx context.Context, projectID string) (*bigquery.Client, error) {
	bqClient, err := client.NewBigqueryClient(ctx, projectID)
	if err != nil {
		return nil, err
	}
	return bqClient, nil
}

func (_ *ChangeStremsWatcherClientImpl) newPubsubClient(ctx context.Context, projectID string) (*pubsub.Client, error) {
	psClient, err := client.NewPubsubClient(ctx, projectID)
	if err != nil {
		return nil, err
	}
	return psClient, nil
}

func (_ *ChangeStremsWatcherClientImpl) newKinesisClient(ctx context.Context) (*kinesis.Client, error) {
	ksClient, err := client.NewKinesisClient(ctx)
	if err != nil {
		return nil, err
	}
	return ksClient, nil
}

func (c *ChangeStremsWatcherClientImpl) watch(ctx context.Context, ops *options.ChangeStreamOptions) (*mongo.ChangeStream, error) {
	cs, err := mongoConnection.Watch(ctx, c.MongoClient, ops)
	if err != nil {
		return nil, err
	}
	return cs, nil
}

// mainに色々と処理を持たせるのが嫌なので、必要な構造体は後入れ
func (c *ChangeStremsWatcherClientImpl) setCsExporter(exporter ChangeStreamsExporterImpl) {
	c.CsExporter = exporter
}

func (c *ChangeStremsWatcherClientImpl) exportChangeStreams(ctx context.Context) error {
	if err := c.CsExporter.exportChangeStreams(ctx); err != nil {
		return err
	}
	return nil
}

func (c *ChangeStremsWatcherImpl) WatchChangeStreams(ctx context.Context) error {
	nowTime, err := common.FetchNowTime()
	if err != nil {
		return err
	}

	pv, err := config.FetchPersistentVolumeDir()
	if err != nil {
		return err
	}

	file := pv + nowTime.Format("2006/01/02/2006-01-02.dat")

	rtByte, err := os.ReadFile(file)

	ops := options.ChangeStream().SetFullDocument(options.UpdateLookup)

	if len(rtByte) == 0 && err == nil {
		c.Log.Info("Failed to get resume token. File is already existed, but resume token is not saved in the file.")
	} else if len(rtByte) == 0 && err != nil {
		c.Log.Info("File saved resume token in is not exists. Get from the current change streams.")
	} else {
		rtStr := string(rtByte)
		var rt interface{} = map[string]string{"_data": strings.TrimRight(rtStr, "\n")}

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
	if err != nil {
		return err
	}

	var (
		bqImpl      interfaceForBigquery.BigqueryImpl
		psImpl        interfaceForPubsub.PubsubImpl
		ksImpl interfaceForKinesisStream.KinesisStreamImpl
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
			psImpl = interfaceForPubsub.PubsubImpl{psClientImpl}
		case KinesisStream:
			ksClient, err := c.Watcher.newKinesisClient(ctx)
			if err != nil {
				return err
			}
			ksClientImpl := &interfaceForKinesisStream.KinesisStreamClientImpl{ksClient}
			ksImpl = interfaceForKinesisStream.KinesisStreamImpl{ksClientImpl}
		default:
			return errors.InternalServerError.Wrap("The export destination is wrong.", fmt.Errorf("You need to set the export destination in the environment variable correctly."))
		}
	}

	rtImpl := interfaceForResumeToken.ResumeTokenImpl{c.Log}

	exporterClient := &changeStreamsExporterClientImpl{cs, bqImpl, psImpl, ksImpl, rtImpl}
	exporter := ChangeStreamsExporterImpl{exporterClient, c.Log}

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
		saveResumeToken(rt string) error
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
		resumeToken   interfaceForResumeToken.ResumeTokenImpl
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
	if err := c.bq.ExportToBigquery(ctx, cs); err != nil {
		return err
	}
	return nil
}

func (c *changeStreamsExporterClientImpl) exportToPubsub(ctx context.Context, cs primitive.M) error {
	if err := c.pubsub.ExportToPubsub(ctx, cs); err != nil {
		return err
	}
	return nil
}

func (c *changeStreamsExporterClientImpl) exportToKinesisStream(ctx context.Context, cs primitive.M) error {
	if err := c.kinesisStream.ExportToKinesisStream(ctx, cs); err != nil {
		return err
	}
	return nil
}

func (c *changeStreamsExporterClientImpl) saveResumeToken(rt string) error {
	if err := c.resumeToken.SaveResumeToken(rt); err != nil {
		return err
	}
	return nil
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
			eDst:= expDstList[i]
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
				default:
					return errors.InternalServerErrorEnvGet.New("The export destination is wrong. You need to set the export destination in the environment variable correctly.")
				}
				return nil
			})
		}

		if err := eg.Wait(); err != nil {
			return err
		}

		csRt := csMap["_id"].(primitive.M)["_data"].(string)

		if err := c.exporter.saveResumeToken(csRt); err != nil {
			return err
		}
	}
	return nil
}
