package errors

import (
	"fmt"
	"runtime"
)

type errType string

const (
	InternalServerError            = errType("500: internal server error")
	InternalServerErrorEnvGet      = errType("500: environment variables get error")
	InternalServerErrorClientGet   = errType("500: client get error")
	InternalServerErrorJsonMarshal = errType("500: json marshal error")
	// mongodb
	InternalServerErrorMongoDbConnect = errType("500: mongodb connect error")
	InternalServerErrorMongoDbOperate = errType("500: mongodb operate error")
	// bigquery
	InternalServerErrorBigqueryInsert = errType("500: bigquery insert error")
	// pubsub
	InternalServerErrorPubSubFind    = errType("500: pubsub find error")
	InternalServerErrorPubSubCreate  = errType("500: pubsub create error")
	InternalServerErrorPubSubPublish = errType("500: pubsub publish error")
	InvalidErrorPubSubOrderingKey    = errType("400: pubsub ordering key error")
	// kinesis stream
	InternalServerErrorKinesisStreamPut = errType("500: kinesis stream put error")
	// local storage file
	InternalServerErrorFilePut = errType("500: file put error")

	//// Storage
	// gcs
	InternalServerErrorGcsCreateNewReader = errType("500: gcs create new reader error")
	InternalServerErrorGcsReader          = errType("500: gcs reader error")
	InternalServerErrorGcsWriteObject     = errType("500: gcs write object error")
	InternalServerErrorGcsNewClient       = errType("500: initialize gcs client error")
	// s3
	InternalServerErrorS3GetObject = errType("500: s3 get object error")
	InternalServerErrorS3PutObject = errType("500: s3 put object error")
	InternalServerErrorS3NewClient = errType("500: initialize s3 client error")
)

func (e errType) New(msg string) error {
	_, file, line, _ := runtime.Caller(1)
	return fmt.Errorf("file: %s, line: %d, errType: %s, orgErrMsg: %s", file, line, e, msg)
}

func (e errType) Wrap(msg string, err error) error {
	_, file, line, _ := runtime.Caller(1)
	return fmt.Errorf("file: %s, line: %d, errType: %s, orgErrMsg: %s, errMsg: %s", file, line, e, msg, err)
}
