package kinesis_stream

import (
	"github.com/cam-inc/mxtransporter/config/constant"
	"os"
)

type KinesisStream struct {
	StreamName          string
	KinesisStreamRegion string
}

func KinesisStreamConfig() KinesisStream {
	var ksCfg KinesisStream
	ksCfg.StreamName = os.Getenv(constant.KINESIS_STREAM_NAME)
	ksCfg.KinesisStreamRegion = os.Getenv(constant.KINESIS_STREAM_REGION)
	return ksCfg
}
