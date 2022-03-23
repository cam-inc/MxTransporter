package resume_token

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/cam-inc/mxtransporter/config"
	"github.com/cam-inc/mxtransporter/config/constant"
	"github.com/cam-inc/mxtransporter/pkg/logger"
	"go.uber.org/zap"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	status := m.Run()

	os.Exit(status)
}

func setEnv(key, value string) error {
	return os.Setenv(key, value)
}

func unsetEnv(key string) error {
	return os.Unsetenv(key)
}

func Test_S3_SaveResumeToken(t *testing.T) {
	var l *zap.SugaredLogger
	logConfig := config.LogConfig()
	l = logger.New(logConfig)

	rt := "00000"

	ctx := context.Background()
	tests := []struct {
		name   string
		runner func(t *testing.T)
	}{
		{
			name: "Failed to fetch RESUME_TOKEN_VOLUME_DIR value.",
			runner: func(t *testing.T) {
				_, err := New(ctx, l)
				if err == nil {
					t.Fatalf("Not behaving as intended.")
				}
			},
		},
		{
			name: "Pass to save resume token in file.",
			runner: func(t *testing.T) {

				currentDir := "mydir"
				bucketName := "dev-mxt-resume-token"
				region := "ap-northeast-1"
				endpoint := ""
				storageType := "s3"

				if err := setEnv(constant.RESUME_TOKEN_VOLUME_TYPE, storageType); err != nil {
					t.Fatalf("Failed to set file RESUME_TOKEN_VOLUME_TYPE environment variables.")
				}

				if err := setEnv(constant.RESUME_TOKEN_VOLUME_DIR, currentDir); err != nil {
					t.Fatalf("Failed to set file RESUME_TOKEN_VOLUME_DIR environment variables.")
				}

				if err := setEnv(constant.MONGODB_COLLECTION, "test"); err != nil {
					t.Fatalf("Failed to set file MONGODB_COLLECTION environment variables.")
				}

				if err := setEnv(constant.RESUME_TOKEN_BUCKET_REGION, region); err != nil {
					t.Fatalf("Failed to set file BUCKET_REGION environment variables.")
				}

				if err := setEnv(constant.RESUME_TOKEN_VOLUME_BUCKET_NAME, bucketName); err != nil {
					t.Fatalf("Failed to set file BUCKET_NAME environment variables.")
				}

				if endpoint != "" {
					if err := setEnv(constant.BUCKET_ENDPOINT, endpoint); err != nil {
						t.Fatalf("Failed to set file BUCKET_ENDPOINT environment variables.")
					}
				}

				if err := setEnv(constant.RESUME_TOKEN_VOLUME_BUCKET_NAME, bucketName); err != nil {
					t.Fatalf("Failed to set file RESUME_TOKEN_VOLUME_BUCKET_NAME environment variables.")
				}

				resumeTokenImpl, err := New(ctx, l)
				if err != nil {
					t.Fatalf("Testing Error, ErrorMessage: %v", err)
				}
				/*
					imp, exists := resumeTokenImpl.(*rtS3)
					if !exists {
						t.Fatalf("Failed resumeTokenImpl convert rtS3.")
					}
					imp.client = &mock{}
				*/
				env := resumeTokenImpl.Env()
				if env == "" {
					t.Fatalf("Failed to get environment variables for resume tokens settings.")
				}
				fmt.Printf("env %s\n", env)

				resumeTokenP := resumeTokenImpl.ReadResumeToken(ctx)
				if resumeTokenP == "" {
					t.Fatal("Failed to read file saved test resume token in.")
				}
				if err := resumeTokenImpl.SaveResumeToken(ctx, rt); err != nil {
					t.Fatalf("Testing Error, ErrorMessage: %v", err)
				}

				resumeToken := resumeTokenImpl.ReadResumeToken(ctx)
				if resumeToken == "" {
					t.Fatal("Failed to read file saved test resume token in.")
				}
				fmt.Printf("resumeToken %s\n", resumeToken)
				envMap := map[string]string{}
				if err := json.Unmarshal([]byte(env), &envMap); err != nil {
					t.Fatal(err)
				}

				if err := unsetEnv(constant.RESUME_TOKEN_VOLUME_TYPE); err != nil {
					t.Fatalf("Failed to unset file RESUME_TOKEN_VOLUME_TYPE environment variables.")
				}
				if err := unsetEnv(constant.RESUME_TOKEN_VOLUME_DIR); err != nil {
					t.Fatalf("Failed to unset file RESUME_TOKEN_VOLUME_DIR environment variables.")
				}
				if err := unsetEnv(constant.MONGODB_COLLECTION); err != nil {
					t.Fatalf("Failed to unset file MONGODB_COLLECTION environment variables.")
				}
				if err := unsetEnv(constant.RESUME_TOKEN_BUCKET_REGION); err != nil {
					t.Fatalf("Failed to unset file RESUME_TOKEN_BUCKET_REGION environment variables.")
				}
				if err := unsetEnv(constant.RESUME_TOKEN_VOLUME_BUCKET_NAME); err != nil {
					t.Fatalf("Failed to unset file RESUME_TOKEN_VOLUME_BUCKET_NAME environment variables.")
				}
				if err := unsetEnv(constant.BUCKET_ENDPOINT); err != nil {
					t.Fatalf("Failed to unset file BUCKET_ENDPOINT environment variables.")
				}
				if err := unsetEnv(constant.RESUME_TOKEN_VOLUME_BUCKET_NAME); err != nil {
					t.Fatalf("Failed to unset file RESUME_TOKEN_VOLUME_BUCKET_NAME environment variables.")
				}

			},
		},
	}

	for _, v := range tests {
		t.Run(v.name, v.runner)
	}
}
