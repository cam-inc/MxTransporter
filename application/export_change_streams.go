package application

import (
	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/pubsub"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/kinesis"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/sync/errgroup"
	"mxtransporter/config"
	interfaceForBigquery "mxtransporter/interfaces/bigquery"
	interfaceKinesisStream "mxtransporter/interfaces/kinesis-stream"
	mongoConnection "mxtransporter/interfaces/mongo"
	interfaceForPubsub "mxtransporter/interfaces/pubsub"
	"mxtransporter/pkg/client"
	"mxtransporter/pkg/errors"
	"mxtransporter/usecases/resume-token"
	"os"
	"strings"
	"time"
)

type Agent string

const (
	BigQuery      Agent = "bigquery"
	CloudPubSub   Agent = "pubsub"
	KinesisStream Agent = "kinesisStream"
)

type (
	bqExpFuncs interface {
		Exec(ctx context.Context, csMap primitive.M, bqClient *bigquery.Client) error
	}
	pubsubExpFuncs interface {
		Exec(ctx context.Context, csMap primitive.M, psClient *pubsub.Client) error
	}
	kinesisStreamExpFuncs interface {
		Exec(ctx context.Context, csMap primitive.M, ksClient *kinesis.Client) error
	}

	bqExpFuncsImpl struct {
	}
	//bqExpFuncsMock struct {
	//}
	pubsubExpFuncsImpl struct {
	}
	//pubsubExpFuncsMock struct {
	//}
	kinesisStreamExpFuncsImpl struct {
	}
	//kinesisStreamExpFuncsMock struct {
	//}

	Exporter struct {
		bq            bqExpFuncs
		pubsub        pubsubExpFuncs
		kinesisStream kinesisStreamExpFuncs
	}
)

func (e *Exporter) BQ(ctx context.Context, csMap primitive.M, bqClient *bigquery.Client) error {
	err := e.bq.Exec(ctx, csMap, bqClient)
	if err != nil {
		return err
	}
	return nil
}
func (e *Exporter) PubSub(ctx context.Context, csMap primitive.M, psClient *pubsub.Client) error {
	err := e.pubsub.Exec(ctx, csMap, psClient)
	if err != nil {
		return err
	}
	return nil
}
func (e *Exporter) KinesisStream(ctx context.Context, csMap primitive.M, ksClient *kinesis.Client) error {
	err := e.kinesisStream.Exec(ctx, csMap, ksClient)
	if err != nil {
		return err
	}
	return nil
}

func (b *bqExpFuncsImpl) Exec(ctx context.Context, csMap primitive.M, bqClient *bigquery.Client) error {
	var bqFunc = &bigqueryClientImpl{
		putRecord: func(ctx context.Context, dataset string, table string, csItems []interfaceForBigquery.ChangeStreamTableSchema, bqClient *bigquery.Client) error {
			return interfaceForBigquery.PutRecord(ctx, dataset, table, csItems, bqClient)
		},
	}

	if err := interfaceForBigquery.NewBigqueryClient(bqFunc).ExportToBigquery(ctx, csMap, bqClient); err != nil {
		return err
	}
	return nil
}

//func (b *bqExpFuncsMock) Exec() {
//}
func (p *pubsubExpFuncsImpl) Exec(ctx context.Context, csMap primitive.M, psClient *pubsub.Client) error {
	var pubsubFunc = &pubsubClientImpl{
		pubsubTopic: func(ctx context.Context, topicID string, psClient *pubsub.Client) error {
			return interfaceForPubsub.PubsubTopic(ctx, topicID, psClient)
		},
		pubsubSubscription: func(ctx context.Context, topicID string, subscriptionID string, psClient *pubsub.Client) error {
			return interfaceForPubsub.PubsubSubscription(ctx, topicID, subscriptionID, psClient)
		},
		publishMessage: func(ctx context.Context, topicID string, csArray []string, psClient *pubsub.Client) error {
			return interfaceForPubsub.PublishMessage(ctx, topicID, csArray, psClient)
		},
	}

	if err := interfaceForPubsub.NewPubsubClient(pubsubFunc).ExportToPubSub(ctx, csMap, psClient); err != nil {
		return err
	}
	return nil
}

//func (b *pubsubExpFuncsMock) Exec() {
//}
func (k *kinesisStreamExpFuncsImpl) Exec(ctx context.Context, csMap primitive.M, ksClient *kinesis.Client) error {
	var kinesisStreamFunc = &kinesisClientImpl{
		putRecord: func(ctx context.Context, streamName string, rt interface{}, csArray []string, ksClient *kinesis.Client) error {
			return interfaceKinesisStream.PutRecord(ctx, streamName, rt, csArray, ksClient)
		},
	}

	if err := interfaceKinesisStream.NewKinesisClient(kinesisStreamFunc).ExportToKinesisStream(ctx, csMap, ksClient); err != nil {
		return err
	}
	return nil
}

//func (b *kinesisStreamExpFuncsMock) Exec() {
//}

func createExporter() *Exporter {
	//if Env == "test" {
	//	return &Exporter{
	//		bq:      &bqClientMock{},
	//		pubsub:  &pubsubClientMock{},
	//		kinesis: &kinesisClientMock{},
	//	}
	//}
	return &Exporter{
		bq:            &bqExpFuncsImpl{},
		pubsub:        &pubsubExpFuncsImpl{},
		kinesisStream: &kinesisStreamExpFuncsImpl{},
	}
}

// export処理(/interface)内で使われているfunctionを定義
// 実際のaws, gcp clientをmock可能にするため
type (
	bigqueryClientImpl struct {
		interfaceForBigquery.BigqueryClientImpl
		putRecord func(ctx context.Context, dataset string, table string, csItems []interfaceForBigquery.ChangeStreamTableSchema, bqClient *bigquery.Client) error
	}
	pubsubClientImpl struct {
		interfaceForPubsub.PubsubClientImple
		pubsubTopic        func(ctx context.Context, topicID string, psClient *pubsub.Client) error
		pubsubSubscription func(ctx context.Context, topicID string, subscriptionID string, psClient *pubsub.Client) error
		publishMessage     func(ctx context.Context, topicID string, csArray []string, psClient *pubsub.Client) error
	}

	kinesisClientImpl struct {
		interfaceKinesisStream.KinesisClientImpl
		putRecord func(ctx context.Context, streamName string, rt interface{}, csArray []string, ksClient *kinesis.Client) error
	}
)

func (b *bigqueryClientImpl) PutRecord(ctx context.Context, dataset string, table string, csItems []interfaceForBigquery.ChangeStreamTableSchema, bqClient *bigquery.Client) error {
	return b.putRecord(ctx, dataset, table, csItems, bqClient)
}
func (p *pubsubClientImpl) PubsubTopic(ctx context.Context, topicID string, psClient *pubsub.Client) error {
	return p.pubsubTopic(ctx, topicID, psClient)
}
func (p *pubsubClientImpl) PubsubSubscription(ctx context.Context, topicID string, subscriptionID string, psClient *pubsub.Client) error {
	return p.pubsubSubscription(ctx, topicID, subscriptionID, psClient)
}
func (p *pubsubClientImpl) PublishMessage(ctx context.Context, topicID string, csArray []string, psClient *pubsub.Client) error {
	return p.publishMessage(ctx, topicID, csArray, psClient)
}
func (k *kinesisClientImpl) PutRecord(ctx context.Context, streamName string, rt interface{}, csArray []string, ksClient *kinesis.Client) error {
	return k.putRecord(ctx, streamName, rt, csArray, ksClient)
}

//=======================================================
// Main function
//=======================================================
func WatchChangeStreams(ctx context.Context, client *mongo.Client) error {
	db, err := mongoConnection.FetchDatabase(ctx, client)
	if err != nil {
		return err
	}
	coll, err := mongoConnection.FetchCollection(ctx, db)
	if err != nil {
		return err
	}

	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		return errors.InternalServerError.Wrap("Failed to load location time.", err)
	}

	nowTime := time.Now().In(jst)

	pv, err := config.FetchPersistentVolumeDir()
	if err != nil {
		return err
	}

	file := pv + nowTime.Format("2006/01/02/2006-01-02.dat")

	rtByte, err := os.ReadFile(file)

	ops := options.ChangeStream().SetFullDocument(options.UpdateLookup)

	if len(rtByte) == 0 && err == nil {
		fmt.Println("Failed to get resume-token. File is already existed, but resume-token is not saved in the file.")
	} else if len(rtByte) == 0 && err != nil {
		fmt.Println("File saved resume-token in is not exists. Get from the current change streams.")
	} else {
		rtStr := string(rtByte)
		var rt interface{} = map[string]string{"_data": strings.TrimRight(rtStr, "\n")}

		ops.SetResumeAfter(rt)
	}

	cs, err := coll.Watch(ctx, mongo.Pipeline{}, ops)
	if err != nil {
		return errors.InternalServerErrorMongoDbOperate.Wrap("Failed to watch mongodb.", err)
	}
	if err := exportChangeStreams(ctx, cs); err != nil {
		return err
	}

	return nil
}

func exportChangeStreams(ctx context.Context, cs *mongo.ChangeStream) error {
	defer cs.Close(ctx)

	exportDestinations, err := config.FetchExportDestination()

	if err != nil {
		return err
	}
	exportDestinationList := strings.Split(exportDestinations, ",")

	exporter := createExporter()

	gcpConfig := config.FetchGcpProject()
	projectID := gcpConfig.ProjectID
	var bqClient *bigquery.Client
	var psClient *pubsub.Client
	var ksClient *kinesis.Client

	for i := 0; i < len(exportDestinationList); i++ {
		exportDestination := exportDestinationList[i]
		var err error
		switch Agent(exportDestination) {
		case BigQuery:
			bqClient, err = client.NewBigqueryClient(ctx, projectID)
			if err != nil {
				return err
			}
		case CloudPubSub:
			psClient, err = client.NewPubsubClient(ctx, projectID)
			if err != nil {
				return err
			}
		case KinesisStream:
			ksClient, err = client.NewKinesisClient(ctx)
			if err != nil {
				return err
			}
		default:
			return errors.InternalServerError.Wrap("The export destination is wrong.", fmt.Errorf("You need to set the export destination in the environment variable correctly."))

		}
	}

	for cs.Next(ctx) {
		var csMap bson.M
		if err := cs.Decode(&csMap); err != nil {
			return errors.InternalServerError.Wrap("Failed to decode change stream.", err)
		}

		csDb := csMap["ns"].(bson.M)["db"]
		csColl := csMap["ns"].(bson.M)["coll"]
		csOpType := csMap["operationType"]
		csClusterTimeInt := csMap["clusterTime"].(primitive.Timestamp)

		fmt.Println(fmt.Sprintf("[INFO] msg: Success to get change-streams, database: %s, collection: %s, operationType: %s, updateTime: %s", csDb, csColl, csOpType, time.Unix(int64(csClusterTimeInt.T), 0)))

		var eg errgroup.Group
		for i := 0; i < len(exportDestinationList); i++ {
			exportDestination := exportDestinationList[i]
			eg.Go(func() error {
				switch Agent(exportDestination) {
				case BigQuery:
					err := exporter.BQ(ctx, csMap, bqClient)
					if err != nil {
						return nil
					}
				case CloudPubSub:
					err := exporter.PubSub(ctx, csMap, psClient)
					if err != nil {
						return nil
					}
				case KinesisStream:
					err := exporter.KinesisStream(ctx, csMap, ksClient)
					if err != nil {
						return nil
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

		if err := resume_token.NewGeneralConfig(pvFunction).SaveResumeToken(csMap["_id"].(primitive.M)); err != nil {
			return err
		}
	}
	return nil
}

// 下は後々、各exporterの呼び出しと形を合わせます
// resume token保存処理のパーツ
//=======================================================
type generalConfig struct {
	config.GeneralConfigIf
	persistentVolume func() (string, error)
}

func (m *generalConfig) FetchPersistentVolumeDir() (string, error) {
	return m.persistentVolume()
}

var pvFunction = &generalConfig{
	persistentVolume: func() (string, error) {
		return config.FetchPersistentVolumeDir()
	},
}

//=======================================================
