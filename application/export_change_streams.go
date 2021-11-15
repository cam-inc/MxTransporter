package application

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/sync/errgroup"
	"mxtransporter/config"
	interfaceForBigquery "mxtransporter/interfaces/bigquery"
	interfaceKinesisStream "mxtransporter/interfaces/kinesis-stream"
	mongoConnection "mxtransporter/interfaces/mongo"
	interfaceForPubSub "mxtransporter/interfaces/pubsub"
	"mxtransporter/pkg/errors"
	"mxtransporter/usecases/resume-token"
	"os"
	"strings"
	"time"
)

type GcpService string
type AwsService string

const (
	Bigquery GcpService = "bigquery"
	PubSub   GcpService = "pubsub"
)

const (
	KinesisStream AwsService = "kinesisStream"
)

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
	nowYear := nowTime.Format("2006")
	nowMonth := nowTime.Format("01")
	nowDay := nowTime.Format("02")

	pv, err := config.FetchPersistentVolumeDir()
	if err != nil {
		return err
	}

	fileName := nowTime.Format("2006-01-02")
	filePath := pv + nowYear + "/" + nowMonth + "/" + nowDay + "/"
	file := filePath + fileName + ".dat"

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
				switch {
				case GcpService(exportDestination) == Bigquery:
					if err := interfaceForBigquery.ExportToBigquery(ctx, csMap, &interfaceForBigquery.BigqueryFuncs{}); err != nil {
						return err
					}
				case GcpService(exportDestination) == PubSub:
					if err := interfaceForPubSub.ExportToPubSub(ctx, csMap, &interfaceForPubSub.PubsubFuncs{}); err != nil {
						return err
					}
				case AwsService(exportDestination) == KinesisStream:
					if err := interfaceKinesisStream.ExportToKinesisStream(ctx, csMap, &interfaceKinesisStream.KinesisFuncs{}); err != nil {
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

		if err := resume_token.NewGeneralConfig(pvFunction).SaveResumeToken(csMap["_id"].(primitive.M)); err != nil {
			return err
		}
	}
	return nil
}

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
