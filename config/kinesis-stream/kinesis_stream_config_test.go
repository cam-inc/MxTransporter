//go:build test
// +build test

package kinesis_stream

import (
	"os"
	"reflect"
	"testing"
)

func Test_KinesisStreamConfig(t *testing.T) {
	t.Run("Check to call the set environment variable.", func(t *testing.T) {
		ksStreamName := "xxx"
		ksRegion := "ap-northeast-1"
		if err := os.Setenv("KINESIS_STREAM_NAME", ksStreamName); err != nil {
			t.Fatalf("Failed to set file KINESIS_STREAM_NAME environment variables.")
		}
		if err := os.Setenv("KINESIS_STREAM_REGION", ksRegion); err != nil {
			t.Fatalf("Failed to set file KINESIS_STREAM_REGION environment variables.")
		}

		ksCfg := KinesisStreamConfig()
		if e, a := ksCfg.StreamName, ksStreamName; !reflect.DeepEqual(e, a) {
			t.Fatal("Environment variable KINESIS_STREAM_NAME is not acquired correctly.")
		}
		if e, a := ksCfg.KinesisStreamRegion, ksRegion; !reflect.DeepEqual(e, a) {
			t.Fatal("Environment variable KINESIS_STREAM_REGION is not acquired correctly.")
		}
	})
}
