package opensearch

import (
	"os"
	"strconv"

	"github.com/cam-inc/mxtransporter/config/constant"
)

type OpenSearch struct {
	EndPoint                      string
	IndexName                     string
	SyncEnabled                   bool
	BulkEnabled                   bool
	BulkFlushBytes                int
	BulkSyncAggregationEnabled    bool
	UseAmazonOpenSearchService    bool
	UseAmazonOpenSearchServerless bool
	AwsRegion                     string
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
	osCfg.EndPoint = os.Getenv(constant.OPENSEARCH_END_POINT)
	osCfg.IndexName = os.Getenv(constant.OPENSEARCH_INDEX_NAME)
	osCfg.SyncEnabled = getBoolEnvWithDefault(constant.OPENSEARCH_SYNC_ENABLED, false)
	osCfg.BulkEnabled = getBoolEnvWithDefault(constant.OPENSEARCH_BULK_ENABLED, false)
	osCfg.BulkFlushBytes = getIntEnvWithDefault(constant.OPENSEARCH_BULK_FLUSH_BYTES, 5e+6)
	osCfg.BulkSyncAggregationEnabled = getBoolEnvWithDefault(constant.OPENSEARCH_BULK_SYNC_AGGREGATION_ENABLED, false)
	osCfg.UseAmazonOpenSearchService = getBoolEnvWithDefault(constant.OPENSEARCH_USE_AMAZON_OPENSEARCH_SERVICE, false)
	osCfg.UseAmazonOpenSearchServerless = getBoolEnvWithDefault(constant.OPENSEARCH_USE_AMAZON_OPENSEARCH_SERVERLESS, false)
	osCfg.AwsRegion = os.Getenv(constant.OPENSEARCH_AWS_REGION)

	return osCfg
}
