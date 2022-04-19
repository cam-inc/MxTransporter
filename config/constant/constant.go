package constant

const (
	BIGQUERY_DATASET = "BIGQUERY_DATASET"
	BIGQUERY_TABLE   = "BIGQUERY_TABLE"

	KINESIS_STREAM_NAME   = "KINESIS_STREAM_NAME"
	KINESIS_STREAM_REGION = "KINESIS_STREAM_REGION"

	MONGODB_HOST       = "MONGODB_HOST"
	MONGODB_DATABASE   = "MONGODB_DATABASE"
	MONGODB_COLLECTION = "MONGODB_COLLECTION"

	RESUME_TOKEN_VOLUME_DIR         = "RESUME_TOKEN_VOLUME_DIR"
	RESUME_TOKEN_VOLUME_TYPE        = "RESUME_TOKEN_VOLUME_TYPE"
	RESUME_TOKEN_VOLUME_BUCKET_NAME = "RESUME_TOKEN_VOLUME_BUCKET_NAME"
	RESUME_TOKEN_FILE_NAME          = "RESUME_TOKEN_FILE_NAME"
	RESUME_TOKEN_BUCKET_REGION      = "RESUME_TOKEN_BUCKET_REGION"
	RESUME_TOKEN_SAVE_INTERVAL_SEC  = "RESUME_TOKEN_SAVE_INTERVAL_SEC"

	EXPORT_DESTINATION                    = "EXPORT_DESTINATION"
	PROJECT_NAME_TO_EXPORT_CHANGE_STREAMS = "PROJECT_NAME_TO_EXPORT_CHANGE_STREAMS"

	TIME_ZONE = "TIME_ZONE"

	LOG_LEVEL            = "LOG_LEVEL"
	LOG_FORMAT           = "LOG_FORMAT"
	LOG_OUTPUT_DIRECTORY = "LOG_OUTPUT_DIRECTORY"
	LOG_OUTPUT_FILE      = "LOG_OUTPUT_FILE"

	FILE_EXPORTER_LOG_TYPE          = "FILE_EXPORTER_LOG_TYPE"
	FILE_EXPORTER_CHANGE_STREAM_KEY = "FILE_EXPORTER_CHANGE_STREAM_KEY"
	FILE_EXPORTER_TIME_KEY          = "FILE_EXPORTER_TIME_KEY"
	FILE_EXPORTER_NAME_KEY          = "FILE_EXPORTER_NAME_KEY"
)
