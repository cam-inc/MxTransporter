package elasticsearch

import (
	"os"
	"strconv"

	"github.com/cam-inc/mxtransporter/config/constant"
)

type Elasticsearch struct {
	ElasticsearchConnectionUrl string
	IndexName                  string
	SyncEnabled                bool
}

func getEnvAsBool(name string, defaultVal bool) bool {
	valStr := os.Getenv(name)
	if val, err := strconv.ParseBool(valStr); err == nil {
		return val
	}
	return defaultVal
}

func ElasticsearchConfig() Elasticsearch {
	var esCfg Elasticsearch
	esCfg.ElasticsearchConnectionUrl = os.Getenv(constant.ELASTICSEARCH_HOST)
	esCfg.IndexName = os.Getenv(constant.ELASTICSEARCH_INDEX_NAME)
	esCfg.SyncEnabled = getEnvAsBool(constant.ELASTICSEARCH_SYNC_ENABLED, false)
	return esCfg
}
