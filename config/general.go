package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"mxtransporter/pkg/errors"
	"os"
)

type GcpProject struct {
	ProjectID string
}

type AwsProfile struct {
	ProfileName string
}

type Region struct {
	Region string
}

func init() {
	m := godotenv.Load()
	if m != nil {
		fmt.Println("[Warning] If this environment is local machine, you have to create .env file, and set env variables with reference to .env.template .")
	}
}

func PersistentVolume() (string, error) {
	pvDir, pvDirExistence := os.LookupEnv("PERSISTENT_VOLUME_DIR")
	if pvDirExistence == false {
		return "", errors.InternalServerErrorEnvGet.New("PERSISTENT_VOLUME_DIR is not existed in environment variables")
	}
	return pvDir, nil
}

func ExportDestination() (string, error) {
	exportDestination, exportDestinationExistence := os.LookupEnv("EXPORT_DESTINATION")
	if exportDestinationExistence == false {
		return "", errors.InternalServerErrorEnvGet.New("EXPORT_DESTINATION is not existed in environment variables")
	}
	return exportDestination, nil
}

func FetchGcpProject() GcpProject {
	var projectConfig GcpProject
	projectConfig.ProjectID = os.Getenv("PROJECT_NAME_TO_EXPORT_CHANGE_STREAMS")
	return projectConfig
}
