package bigquery

import (
	"os"
)

type Bigquery struct {
	DataSet string
	Table   string
}

func BigqueryConfig() Bigquery {
	var bigqueryConfig Bigquery
	bigqueryConfig.DataSet = os.Getenv("BIGQUERY_DATASET")
	bigqueryConfig.Table = os.Getenv("BIGQUERY_TABLE")
	return bigqueryConfig
}
