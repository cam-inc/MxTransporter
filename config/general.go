package config

import (
	"github.com/cam-inc/mxtransporter/pkg/errors"
	"github.com/cam-inc/mxtransporter/pkg/logger"
	"github.com/joho/godotenv"
	"os"
)

func init() {
	// for runing locally
	godotenv.Load()
}

func FetchPersistentVolumeDir() (string, error) {
	// Required parameter uses lookupEnv ()
	pvDir, pvDirExistence := os.LookupEnv("PERSISTENT_VOLUME_DIR")
	if !pvDirExistence {
		return "", errors.InternalServerErrorEnvGet.New("PERSISTENT_VOLUME_DIR is not existed in environment variables")
	}
	return pvDir, nil
}

func FetchExportDestination() (string, error) {
	expDst, expDstExistence := os.LookupEnv("EXPORT_DESTINATION")
	if !expDstExistence {
		return "", errors.InternalServerErrorEnvGet.New("EXPORT_DESTINATION is not existed in environment variables")
	}
	return expDst, nil
}

func FetchGcpProject() (string, error) {
	// LookupEnv() is used because error judgment is required for error handling of the caller.
	projectID, projectIDExistence := os.LookupEnv("PROJECT_NAME_TO_EXPORT_CHANGE_STREAMS")
	if !projectIDExistence {
		return "", errors.InternalServerErrorEnvGet.New("PROJECT_NAME_TO_EXPORT_CHANGE_STREAMS is not existed in environment variables")
	}
	return projectID, nil
}

func FetchTimeZone() (string, error) {
	tz, tzExistence := os.LookupEnv("TIME_ZONE")
	if !tzExistence {
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
