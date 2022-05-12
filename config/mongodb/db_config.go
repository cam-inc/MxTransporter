package mongodb

import (
	"github.com/cam-inc/mxtransporter/config/constant"
	"os"
)

type Mongo struct {
	MongoDbConnectionUrl             string
	MongoDbDatabase                  string
	MongoDbCollection                string
	MongoWatchPipelineExcludeCsFulldocumentField string
}

func MongoConfig() Mongo {
	var mCfg Mongo
	mCfg.MongoDbConnectionUrl = os.Getenv(constant.MONGODB_HOST)
	mCfg.MongoDbDatabase = os.Getenv(constant.MONGODB_DATABASE)
	mCfg.MongoDbCollection = os.Getenv(constant.MONGODB_COLLECTION)
	mCfg.MongoWatchPipelineExcludeCsFulldocumentField = os.Getenv(constant.MONGO_WATCH_PIPELINE_EXCLUDE_CS_FULLDOCUMENT_FIELD)
	return mCfg
}
