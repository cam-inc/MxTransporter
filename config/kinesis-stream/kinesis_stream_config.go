package kinesis_stream

import (
	"os"
)

type KinesisStream struct {
	StreamName string
}

func KinesisStreamConfig() KinesisStream {
	var kinesisStreamConfig KinesisStream
	kinesisStreamConfig.StreamName = os.Getenv("KINESIS_STREAM_NAME")
	return kinesisStreamConfig
}