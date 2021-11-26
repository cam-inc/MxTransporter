package errors

import (
	"fmt"
	"runtime"
)

type ErrorType string

const (
	InternalServerError            = ErrorType("500: internal server error")
	InternalServerErrorEnvGet      = ErrorType("500: environment variables get error")
	InternalServerErrorClientGet   = ErrorType("500: client get error")
	InternalServerErrorJsonMarshal = ErrorType("500: json marshal error")
	// mongodb
	InternalServerErrorMongoDbConnect = ErrorType("500: mongodb connect error")
	InternalServerErrorMongoDbOperate = ErrorType("500: mongodb operate error")
	// bigquery
	InternalServerErrorBigqueryInsert = ErrorType("500: bigquery insert error")
	// pubsub
	InternalServerErrorPubSubFind   = ErrorType("500: pubsub find error")
	InternalServerErrorPubSubCreate = ErrorType("500: pubsub create error")
	// kinesis stream
	InternalServerErrorKinesisStreamPut = ErrorType("500: kinesis stream put error")
)

func (et ErrorType) New(msg string) error {
	_, file, line, _ := runtime.Caller(1)
	return fmt.Errorf("[ERROR] file: %s, line: %d, errorType: %s, originalErrorMessage: %s", file, line, et, msg)
}

func (et ErrorType) Wrap(msg string, err error) error {
	_, file, line, _ := runtime.Caller(1)
	return fmt.Errorf("[ERROR] file: %s, line: %d, errorType: %s, originalErrorMessage: %s, errorMessage: %s", file, line, et, msg, err)
}
