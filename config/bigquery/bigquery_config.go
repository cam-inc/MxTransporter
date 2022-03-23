package bigquery

import (
	"github.com/cam-inc/mxtransporter/config/constant"
	"os"
)

type Bigquery struct {
	DataSet string
	Table   string
}

func BigqueryConfig() Bigquery {
	var bqCfg Bigquery
	bqCfg.DataSet = os.Getenv(constant.BIGQUERY_DATASET)
	bqCfg.Table = os.Getenv(constant.BIGQUERY_TABLE)
	return bqCfg
}
