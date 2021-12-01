package bigquery

import (
	"os"
)

type Bigquery struct {
	DataSet string
	Table   string
}

func BigqueryConfig() Bigquery {
	var bqCfg Bigquery
	bqCfg.DataSet = os.Getenv("BIGQUERY_DATASET")
	bqCfg.Table = os.Getenv("BIGQUERY_TABLE")
	return bqCfg
}
