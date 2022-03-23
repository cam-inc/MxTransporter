package resume_token

import (
	"github.com/cam-inc/mxtransporter/config/constant"
	"os"
	"strconv"
)

type (
	ResumeToken struct {
		VolumeType      string
		BucketName      string
		Region          string
		Path            string
		SaveIntervalSec int
	}
)

func ResumeTokenConfig() ResumeToken {
	var config ResumeToken
	config.Path = os.Getenv(constant.RESUME_TOKEN_VOLUME_DIR)
	config.VolumeType = os.Getenv(constant.RESUME_TOKEN_VOLUME_TYPE)
	config.BucketName = os.Getenv(constant.RESUME_TOKEN_VOLUME_BUCKET_NAME)
	config.Region = os.Getenv(constant.RESUME_TOKEN_BUCKET_REGION)
	interval := os.Getenv(constant.RESUME_TOKEN_SAVE_INTERVAL_SEC)
	intervalSec, _ := strconv.Atoi(interval)
	config.SaveIntervalSec = intervalSec
	return config
}
