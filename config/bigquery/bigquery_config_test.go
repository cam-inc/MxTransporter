//go:build test
// +build test

package bigquery

import (
	"os"
	"reflect"
	"testing"
)

func Test_BigqueryConfig(t *testing.T) {
	t.Run("Check to call the set environment variable.", func(t *testing.T) {
		bqDataset := "xxx"
		bqTable := "yyy"
		if err := os.Setenv("BIGQUERY_DATASET", bqDataset); err != nil {
			t.Fatalf("Failed to set file BIGQUERY_DATASET environment variables.")
		}
		if err := os.Setenv("BIGQUERY_TABLE", bqTable); err != nil {
			t.Fatalf("Failed to set file BIGQUERY_TABLE environment variables.")
		}

		bqCfg := BigqueryConfig()
		if e, a := bqCfg.DataSet, bqDataset; !reflect.DeepEqual(e, a) {
			t.Fatal("Environment variable BIGQUERY_DATASET is not acquired correctly.")
		}
		if e, a := bqCfg.Table, bqTable; !reflect.DeepEqual(e, a) {
			t.Fatal("Environment variable BIGQUERY_TABLE is not acquired correctly.")
		}
	})
}
