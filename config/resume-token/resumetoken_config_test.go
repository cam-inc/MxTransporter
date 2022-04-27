//go:build test
// +build test

package resume_token

import (
	"github.com/cam-inc/mxtransporter/config/constant"
	"os"
	"reflect"
	"testing"
)

func setEnv(key, value string) error {
	return os.Setenv(key, value)
}

func Test_ResumeTokenConfig(t *testing.T) {
	t.Run("Check to call the set environment variable.", func(t *testing.T) {

		if err := setEnv(constant.RESUME_TOKEN_VOLUME_TYPE, "file"); err != nil {
			t.Fatalf("Failed to set file RESUME_TOKEN_VOLUME_TYPE environment variables.")
		}

		if err := setEnv(constant.RESUME_TOKEN_VOLUME_DIR, "path"); err != nil {
			t.Fatalf("Failed to set file RESUME_TOKEN_VOLUME_DIR environment variables.")
		}

		if err := setEnv(constant.MONGODB_COLLECTION, "test"); err != nil {
			t.Fatalf("Failed to set file MONGODB_COLLECTION environment variables.")
		}

		if err := setEnv(constant.RESUME_TOKEN_BUCKET_REGION, "asia-northeast1"); err != nil {
			t.Fatalf("Failed to set file BUCKET_REGION environment variables.")
		}

		if err := setEnv(constant.RESUME_TOKEN_VOLUME_BUCKET_NAME, "my-bucket"); err != nil {
			t.Fatalf("Failed to set file BUCKET_NAME environment variables.")
		}

		if err := setEnv(constant.RESUME_TOKEN_SAVE_INTERVAL_SEC, "10"); err != nil {
			t.Fatalf("Failed to set file RESUME_TOKEN_SAVE_INTERVAL_SEC environment variables.")
		}

		cfg := ResumeTokenConfig()
		if e, a := cfg.VolumeType, "file"; !reflect.DeepEqual(e, a) {
			t.Fatal("Environment variable VolumeType is not acquired correctly.")
		}
		if e, a := cfg.Path, "path"; !reflect.DeepEqual(e, a) {
			t.Fatal("Environment variable Path is not acquired correctly.")
		}
		if e, a := cfg.BucketName, "my-bucket"; !reflect.DeepEqual(e, a) {
			t.Fatal("Environment variable BucketName is not acquired correctly.")
		}
		if e, a := cfg.SaveIntervalSec, 10; !reflect.DeepEqual(e, a) {
			t.Fatal("Environment variable SaveIntervalSec is not acquired correctly.")
		}
	})
}
