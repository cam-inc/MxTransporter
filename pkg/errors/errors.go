package errors

import (
	"fmt"
	"runtime"
)

type ErrorType string

const (
	InternalServerError = ErrorType("General Error")
	// mongodb
	InternalServerErrorMongoDb = ErrorType("MongoDB Error")
	// bigquery
	InternalServerErrorBigquery = ErrorType("BigQuery Error")
	// pubsub
	InternalServerErrorPubSub = ErrorType("PubSub Error")
	// kinesis stream
	InternalServerErrorKinesisStream = ErrorType("Kinesis Stream Error")
	// resume-token
	InternalServerErrorResumeToken = ErrorType("Resume Token Error")
)

func (et ErrorType) New(msg string) error {
	_, file, line, _ := runtime.Caller(1)
	return fmt.Errorf("[ERROR] file: %s, line: %d, errorType: %s, originalErrorMessage: %s", file, line, et, msg)
}

func (et ErrorType) Wrap(msg string, err error) error {
	_, file, line, _ := runtime.Caller(1)
	return fmt.Errorf("[ERROR] file: %s, line: %d, errorType: %s, originalErrorMessage: %s, errorMessage: %s", file, line, et, msg, err)
}