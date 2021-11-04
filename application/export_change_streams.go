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
	Bigquery     	GcpService = "bigquery"
	PubSub			GcpService = "pubsub"
)

const (
	KinesisStream 	AwsService = "kinesisStream"
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

	pv, err := config.PersistentVolume()
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

	exportDestinations, err := config.ExportDestination()

	if err != nil {
		return err
	}
	exportDestinationList := strings.Split(exportDestinations, ",")

	var bigqueryClient *bigquery.Client
	var pubSubClient *pubsub.Client
	var kinesisClient *kinesis.Client

	for i := 0; i < len(exportDestinationList); i++ {
		exportDestination := exportDestinationList[i]

		var err error

		if GcpService(exportDestination) == Bigquery {
			gcpProjectID := config.FetchGcpProject().ProjectID
			bigqueryClient, err = NewBigqueryClient(ctx, gcpProjectID)
			if err != nil {
				return err
			}
		} else if GcpService(exportDestination) == PubSub {
			gcpProjectID := config.FetchGcpProject().ProjectID
			pubSubClient, err = NewPubSubClient(ctx, gcpProjectID)
			if err != nil {
				return err
			}
		} else if AwsService(exportDestination) == KinesisStream {
			kinesisClient, err = NewKinesisClient(ctx)
			if err != nil {
				return err
			}
		}  else {
			return errors.InternalServerErrorEnvGet.New("The export destination is wrong. You need to set the export destination in the environment variable correctly.")
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
				if GcpService(exportDestination) == Bigquery {
					if err := exportToBigquery(ctx, csMap, bigqueryClient); err != nil {
						return err
					}
				} else if GcpService(exportDestination) == PubSub {
					if err := exportToPubSub(ctx, csMap, pubSubClient); err != nil {
						return err
					}
				} else if AwsService(exportDestination) == KinesisStream {
					if err := exportToKinesisStream(ctx, csMap, kinesisClient); err != nil {
						return err
					}
				}  else {
					return errors.InternalServerErrorEnvGet.New("The export destination is wrong. You need to set the export destination in the environment variable correctly.")
				}
				return nil
			})
		}

		if err := eg.Wait(); err != nil {
			return err
		}

		if err := resume_token.SaveResumeToken(csMap["_id"].(primitive.M)); err != nil {
			return err
		}
	}
	return nil
}

func exportToBigquery(ctx context.Context, csMap primitive.M, bigqueryClient *bigquery.Client) error {
	if err := interfaceForBigquery.ExportToBigquery(ctx, csMap, bigqueryClient); err != nil {
		return err
	}
	return nil
}

func exportToPubSub(ctx context.Context, csMap primitive.M, pubSubClient *pubsub.Client) error {
	if err := interfaceForPubSub.ExportToPubSub(ctx, csMap, pubSubClient); err != nil {
		return err
	}
	return nil
}

func exportToKinesisStream(ctx context.Context, csMap primitive.M, kinesisClient *kinesis.Client) error {
	if err := interfaceKinesisStream.ExportToKinesisStream(ctx, csMap, kinesisClient); err != nil {
		return err
	}
	return nil
}