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
	InternalServerErrorPubSubFind   = errType("500: pubsub find error")
	InternalServerErrorPubSubCreate = errType("500: pubsub create error")
	// kinesis stream
	InternalServerErrorKinesisStreamPut = errType("500: kinesis stream put error")
)

func (e errType) New(msg string) error {
	_, file, line, _ := runtime.Caller(1)
	return fmt.Errorf("file: %s, line: %d, errType: %s, orgErrMsg: %s", file, line, e, msg)
}

func (e errType) Wrap(msg string, err error) error {
	_, file, line, _ := runtime.Caller(1)
	return fmt.Errorf("file: %s, line: %d, errType: %s, orgErrMsg: %s, errMsg: %s", file, line, e, msg, err)
}
