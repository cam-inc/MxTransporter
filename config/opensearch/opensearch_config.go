package opensearch

import (
	"os"
	"strconv"

	"github.com/cam-inc/mxtransporter/config/constant"
)

type OpenSearch struct {
	OpenSearchConnectionUrl    string
	IndexName                  string
	SyncEnabled                bool
	BulkEnabled                bool
	BulkFlushBytes             int
	BulkFlushIntervalSeconds   int
	BulkSyncAggregationEnabled bool
}

func getBoolEnvWithDefault(key string, defaultVal bool) bool {
	valStr := os.Getenv(key)
	if valStr == "" {
		return defaultVal
	}

	val, err := strconv.ParseBool(valStr)
	if err != nil {
		return defaultVal
	}

	return val
}

func getIntEnvWithDefault(key string, defaultVal int) int {
	valStr := os.Getenv(key)
	if valStr == "" {
		return defaultVal
	}

	value, err := strconv.Atoi(valStr)
	if err != nil {
		return defaultVal
	}

	return value
}

func OpenSearchConfig() OpenSearch {
	var osCfg OpenSearch
	osCfg.OpenSearchConnectionUrl = os.Getenv(constant.OPENSEARCH_HOST)
	osCfg.IndexName = os.Getenv(constant.OPENSEARCH_INDEX_NAME)
	osCfg.SyncEnabled = getBoolEnvWithDefault(constant.OPENSEARCH_SYNC_ENABLED, false)
	osCfg.BulkEnabled = getBoolEnvWithDefault(constant.OPENSEARCH_BULK_ENABLED, false)
	osCfg.BulkFlushBytes = getIntEnvWithDefault(constant.OPENSEARCH_BULK_FLUSH_BYTES, 5e+6)
	osCfg.BulkFlushIntervalSeconds = getIntEnvWithDefault(constant.OPENSEARCH_BULK_FLUSH_INTERVAL_SECONDS, 30)
	osCfg.BulkSyncAggregationEnabled = getBoolEnvWithDefault(constant.OPENSEARCH_BULK_SYNC_AGGREGATION_ENABLED, false)
	return osCfg
}
