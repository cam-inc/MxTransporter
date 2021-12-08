package kinesis_stream

import (
	"os"
)

type KinesisStream struct {
	StreamName          string
	KinesisStreamRegion string
}

func KinesisStreamConfig() KinesisStream {
	var ksCfg KinesisStream
	ksCfg.StreamName = os.Getenv("KINESIS_STREAM_NAME")
	ksCfg.KinesisStreamRegion = os.Getenv("KINESIS_STREAM_REGION")
	return ksCfg
}
