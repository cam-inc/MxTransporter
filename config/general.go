package config

import (
	"fmt"
	"github.com/cam-inc/mxtransporter/config/constant"
	iff "github.com/cam-inc/mxtransporter/interfaces/file"
	"github.com/cam-inc/mxtransporter/pkg/errors"
	"github.com/cam-inc/mxtransporter/pkg/logger"
	"github.com/joho/godotenv"
	"os"
	"strconv"
	"strings"
)

func init() {
	// for runing locally
	godotenv.Load()
}

func FetchResumeTokenFileName() (string, error) {
	rtFileName, exists := os.LookupEnv(constant.RESUME_TOKEN_FILE_NAME)
	if !exists {
		// default value -> use watching collection name.
		colName, exists := os.LookupEnv(constant.MONGODB_COLLECTION)
		if !exists {
			return "", errors.InternalServerErrorEnvGet.New("MONGODB_COLLECTION is not existed in environment variables")
		}
		rtFileName = fmt.Sprintf("%s.dat", colName)
	}
	return rtFileName, nil
}

func FetchExportDestination() (string, error) {
	expDst, expDstExistence := os.LookupEnv(constant.EXPORT_DESTINATION)
	if !expDstExistence {
		return "", errors.InternalServerErrorEnvGet.New("EXPORT_DESTINATION is not existed in environment variables")
	}
	return expDst, nil
}

func FetchGcpProject() (string, error) {
	// LookupEnv() is used because error judgment is required for error handling of the caller.
	projectID, projectIDExistence := os.LookupEnv(constant.PROJECT_NAME_TO_EXPORT_CHANGE_STREAMS)
	if !projectIDExistence {
		return "", errors.InternalServerErrorEnvGet.New("PROJECT_NAME_TO_EXPORT_CHANGE_STREAMS is not existed in environment variables")
	}
	return projectID, nil
}

func FetchResumeTokenUnusedMode() (bool, error) {
	rtUnusedModeFlag := os.Getenv(constant.RESUME_TOKEN_UNUSED_MODE)

	if rtUnusedModeFlag == "" || strings.EqualFold(rtUnusedModeFlag, "false") {
		return false, nil
	} else if strings.EqualFold(rtUnusedModeFlag, "true") {
		return true, nil
	} else {
		return false, errors.InternalServerError.New("The environment variable RESUME_TOKEN_UNUSED_MODE is not set to the proper value.")
	}
}

func FetchTimeZone() (string, error) {
	tz, tzExistence := os.LookupEnv(constant.TIME_ZONE)
	if !tzExistence {
		return "", errors.InternalServerErrorEnvGet.New("TIME_ZONE is not existed in environment variables")
	}
	return tz, nil
}

func LogConfig() logger.Log {
	var l logger.Log
	l.Level = os.Getenv(constant.LOG_LEVEL)
	l.Format = os.Getenv(constant.LOG_FORMAT)
	l.OutputDirectory = os.Getenv(constant.LOG_OUTPUT_DIRECTORY)
	l.OutputFile = os.Getenv(constant.LOG_OUTPUT_FILE)
	return l
}

func FileExportConfig() *iff.ExporterConfig {
	cfg := &iff.ExporterConfig{}
	cfg.WriterConfig.Writer = os.Getenv(constant.FILE_EXPORTER_WRITER)
	cfg.WriterConfig.MaxMegaBytes, _ = strconv.Atoi(os.Getenv(constant.FILE_EXPORTER_WRITER_MAX_MEGABYTES))
	cfg.WriterConfig.MaxDays, _ = strconv.Atoi(os.Getenv(constant.FILE_EXPORTER_WRITER_MAX_DAYS))
	cfg.WriterConfig.MaxFileBackups, _ = strconv.Atoi(os.Getenv(constant.FILE_EXPORTER_WRITER_MAX_BACKUPS))
	cfg.LogType = os.Getenv(constant.FILE_EXPORTER_LOG_TYPE)
	cfg.ChangeStreamKey = os.Getenv(constant.FILE_EXPORTER_CHANGE_STREAM_KEY)
	cfg.NameKey = os.Getenv(constant.FILE_EXPORTER_NAME_KEY)
	cfg.TimeKey = os.Getenv(constant.FILE_EXPORTER_TIME_KEY)
	return cfg
}
