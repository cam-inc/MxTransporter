package config

import (
	"mxtransporter/pkg/errors"
	"fmt"
	"github.com/joho/godotenv"
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
		return "", errors.InternalServerError.New("PERSISTENT_VOLUME_DIR is not existed in environment variables")
	}
	return pvDir, nil
}

func ExportDestination() (string, error) {
	exportDestination, exportDestinationExistence := os.LookupEnv("EXPORT_DESTINATION")
	if exportDestinationExistence == false {
		return "", errors.InternalServerError.New("EXPORT_DESTINATION is not existed in environment variables")
	}
	return exportDestination, nil
}

func FetchGcpProject() GcpProject {
	var projectConfig GcpProject
	projectConfig.ProjectID = os.Getenv("PROJECT_NAME_TO_EXPORT_CHANGE_STREAMS")
	return projectConfig
}

func FetchAwsProfile() AwsProfile {
	var profileConfig AwsProfile
	profileConfig.ProfileName = os.Getenv("PROFILE_NAME_TO_EXPORT_CHANGE_STREAMS")
	return profileConfig
}

func FetchRegion() Region {
	var regionConfig Region
	regionConfig.Region = os.Getenv("AWS_REGION")
	return regionConfig
}
