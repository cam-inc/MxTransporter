package kinesis_stream

import (
	"os"
)

type KinesisStream struct {
	StreamName 	string
	KinesisStreamRegion		string
}

func KinesisStreamConfig() KinesisStream {
	var kinesisStreamConfig KinesisStream
	kinesisStreamConfig.StreamName = os.Getenv("KINESIS_STREAM_NAME")
	kinesisStreamConfig.KinesisStreamRegion = os.Getenv("KINESIS_STREAM_REGION")
	return kinesisStreamConfig
}