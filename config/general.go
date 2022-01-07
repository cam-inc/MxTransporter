package config

import (
	"github.com/joho/godotenv"
	"mxtransporter/pkg/errors"
	"mxtransporter/pkg/logger"
	"os"
)

func init() {
	// for runing locally
	godotenv.Load()
}

// Required parameter uses lookupEnv ()
func FetchPersistentVolumeDir() (string, error) {
	pvDir, pvDirExistence := os.LookupEnv("PERSISTENT_VOLUME_DIR")
	if pvDirExistence == false {
		return "", errors.InternalServerErrorEnvGet.New("PERSISTENT_VOLUME_DIR is not existed in environment variables")
	}
	return pvDir, nil
}

func FetchExportDestination() (string, error) {
	expDst, expDstExistence := os.LookupEnv("EXPORT_DESTINATION")
	if expDstExistence == false {
		return "", errors.InternalServerErrorEnvGet.New("EXPORT_DESTINATION is not existed in environment variables")
	}
	return expDst, nil
}

// LookupEnv() is used because error judgment is required for error handling of the caller.
func FetchGcpProject() (string, error) {
	projectID, projectIDExistence := os.LookupEnv("PROJECT_NAME_TO_EXPORT_CHANGE_STREAMS")
	if projectIDExistence == false {
		return "", errors.InternalServerErrorEnvGet.New("PROJECT_NAME_TO_EXPORT_CHANGE_STREAMS is not existed in environment variables")
	}
	return projectID, nil
}

func FetchTimeZone() (string, error) {
	tz, tzExistence := os.LookupEnv("TIME_ZONE")
	if tzExistence == false {
		return "", errors.InternalServerErrorEnvGet.New("TIME_ZONE is not existed in environment variables")
	}
	return tz, nil
}

func LogConfig() logger.Log {
	var l logger.Log
	l.Level = os.Getenv("LOG_LEVEL")
	l.Format = os.Getenv("LOG_FORMAT")
	l.OutputDirectory = os.Getenv("LOG_OUTPUT_DIRECTORY")
	l.OutputFile = os.Getenv("LOG_OUTPUT_FILE")
	return l
}
