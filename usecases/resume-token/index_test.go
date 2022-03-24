//go:build test
// +build test

package resume_token

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/cam-inc/mxtransporter/config"
	"github.com/cam-inc/mxtransporter/config/constant"
	"github.com/cam-inc/mxtransporter/pkg/logger"
	mocks "github.com/cam-inc/mxtransporter/usecases/resume-token/mock"
	"github.com/golang/mock/gomock"
	"go.uber.org/zap"
	"os"
	"strings"
	"testing"
	"time"
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

func Test_SaveResumeToken(t *testing.T) {
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
			name: "Pass to read resume token",
			runner: func(t *testing.T) {
				currentDir := "mydir"
				bucketName := "mxt-resume-token-test"
				region := "asia-northeast1"
				storageType := "file"
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

				if err := setEnv(constant.RESUME_TOKEN_VOLUME_BUCKET_NAME, bucketName); err != nil {
					t.Fatalf("Failed to set file RESUME_TOKEN_VOLUME_BUCKET_NAME environment variables.")
				}
				i, err := New(ctx, l)
				if err != nil {
					t.Fatalf("Testing Error, ErrorMessage: %v", err)
				}

				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				cli := mocks.NewMockStorageClient(ctrl)
				resumeToken, exists := i.(*resumeTokenImpl)
				if !exists {
					t.Fatalf("Testing Error. convert interaface to struct failed.")
				}
				resumeToken.client = cli

				ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
				defer cancel()

				cli.EXPECT().
					GetObject(ctx, fmt.Sprintf("%s/%s", resumeToken.volumePath, resumeToken.tokenFileName)).
					Return([]byte(rt), nil).AnyTimes()

				env := resumeToken.Env()
				if env == "" {
					t.Fatalf("Failed to get environment variables for resume tokens settings.")
				}
				fmt.Printf("env %s\n", env)

				token := resumeToken.ReadResumeToken(ctx)
				if token == "" {
					t.Fatal("Failed to read file saved test resume token in.")
				}
				if token != rt {
					t.Fatal("token value mismatch")
				}
				fmt.Printf("resumeToken %s\n", token)

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

				if err := unsetEnv(constant.RESUME_TOKEN_VOLUME_BUCKET_NAME); err != nil {
					t.Fatalf("Failed to unset file RESUME_TOKEN_VOLUME_BUCKET_NAME environment variables.")
				}

			},
		},
		{name: "Error to read resume token",
			runner: func(t *testing.T) {
				currentDir := "mydir"
				bucketName := "mxt-resume-token-test"
				region := "asia-northeast1"
				storageType := "file"
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

				if err := setEnv(constant.RESUME_TOKEN_VOLUME_BUCKET_NAME, bucketName); err != nil {
					t.Fatalf("Failed to set file RESUME_TOKEN_VOLUME_BUCKET_NAME environment variables.")
				}
				i, err := New(ctx, l)
				if err != nil {
					t.Fatalf("Testing Error, ErrorMessage: %v", err)
				}

				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				cli := mocks.NewMockStorageClient(ctrl)
				resumeToken, exists := i.(*resumeTokenImpl)
				if !exists {
					t.Fatalf("Testing Error. convert interaface to struct failed.")
				}
				resumeToken.client = cli

				ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
				defer cancel()

				cli.EXPECT().
					GetObject(ctx, fmt.Sprintf("%s/%s", resumeToken.volumePath, resumeToken.tokenFileName)).
					Return(nil, fmt.Errorf("storage error")).AnyTimes()

				env := resumeToken.Env()
				if env == "" {
					t.Fatalf("Failed to get environment variables for resume tokens settings.")
				}
				fmt.Printf("env %s\n", env)

				token := resumeToken.ReadResumeToken(ctx)
				if token != "" {
					t.Fatal("Failed to read-resume-token error case.")
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

				if err := unsetEnv(constant.RESUME_TOKEN_VOLUME_BUCKET_NAME); err != nil {
					t.Fatalf("Failed to unset file RESUME_TOKEN_VOLUME_BUCKET_NAME environment variables.")
				}

			},
		},
		{
			name: "Pass to save resume token",
			runner: func(t *testing.T) {

				currentDir := "mydir"
				bucketName := "mxt-resume-token-test"
				region := "asia-northeast1"
				storageType := "file"

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

				if err := setEnv(constant.RESUME_TOKEN_VOLUME_BUCKET_NAME, bucketName); err != nil {
					t.Fatalf("Failed to set file RESUME_TOKEN_VOLUME_BUCKET_NAME environment variables.")
				}

				i, err := New(ctx, l)
				if err != nil {
					t.Fatalf("Testing Error, ErrorMessage: %v", err)
				}

				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				cli := mocks.NewMockStorageClient(ctrl)
				resumeToken, exists := i.(*resumeTokenImpl)
				if !exists {
					t.Fatalf("Testing Error. convert interaface to struct failed.")
				}
				resumeToken.client = cli

				ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
				defer cancel()

				cli.EXPECT().
					PutObject(ctx, fmt.Sprintf("%s/%s", resumeToken.volumePath, resumeToken.tokenFileName), rt).
					Return(nil).AnyTimes()

				cli.EXPECT().
					GetObject(ctx, fmt.Sprintf("%s/%s", resumeToken.volumePath, resumeToken.tokenFileName)).
					Return([]byte(rt), nil).AnyTimes()

				env := resumeToken.Env()
				if env == "" {
					t.Fatalf("Failed to get environment variables for resume tokens settings.")
				}
				fmt.Printf("env %s\n", env)

				token := resumeToken.ReadResumeToken(ctx)
				if token == "" {
					t.Fatal("Failed to read file saved test resume token in.")
				}
				if err := resumeToken.SaveResumeToken(ctx, rt); err != nil {
					t.Fatalf("Testing Error, ErrorMessage: %v", err)
				}

				token = resumeToken.ReadResumeToken(ctx)
				if token == "" {
					t.Fatal("Failed to read file saved test resume token in.")
				}
				fmt.Printf("resumeToken %s\n", token)
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

				if err := unsetEnv(constant.RESUME_TOKEN_VOLUME_BUCKET_NAME); err != nil {
					t.Fatalf("Failed to unset file RESUME_TOKEN_VOLUME_BUCKET_NAME environment variables.")
				}

			},
		},
		{
			name: "Error to save resume token",
			runner: func(t *testing.T) {
				currentDir := "mydir"
				bucketName := "mxt-resume-token-test"
				region := "asia-northeast1"
				storageType := "file"

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

				if err := setEnv(constant.RESUME_TOKEN_VOLUME_BUCKET_NAME, bucketName); err != nil {
					t.Fatalf("Failed to set file RESUME_TOKEN_VOLUME_BUCKET_NAME environment variables.")
				}

				i, err := New(ctx, l)
				if err != nil {
					t.Fatalf("Testing Error, ErrorMessage: %v", err)
				}

				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				cli := mocks.NewMockStorageClient(ctrl)
				resumeToken, exists := i.(*resumeTokenImpl)
				if !exists {
					t.Fatalf("Testing Error. convert interaface to struct failed.")
				}
				resumeToken.client = cli

				ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
				defer cancel()

				expectErr := fmt.Errorf("storage error")

				cli.EXPECT().
					PutObject(ctx, fmt.Sprintf("%s/%s", resumeToken.volumePath, resumeToken.tokenFileName), rt).
					Return(expectErr).AnyTimes()

				env := resumeToken.Env()
				if env == "" {
					t.Fatalf("Failed to get environment variables for resume tokens settings.")
				}
				fmt.Printf("env %s\n", env)

				if err := resumeToken.SaveResumeToken(ctx, rt); err == nil || !strings.Contains(err.Error(), expectErr.Error()) {
					t.Fatalf("Testing Error, ErrorMessage: %v <-> %v", err, expectErr)
				}

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

				if err := unsetEnv(constant.RESUME_TOKEN_VOLUME_BUCKET_NAME); err != nil {
					t.Fatalf("Failed to unset file RESUME_TOKEN_VOLUME_BUCKET_NAME environment variables.")
				}
			},
		},
		{
			name: "Skip to save resume token",
			runner: func(t *testing.T) {
				currentDir := "mydir"
				bucketName := "mxt-resume-token-test"
				region := "asia-northeast1"
				storageType := "file"
				intavalSec := "10"

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

				if err := setEnv(constant.RESUME_TOKEN_VOLUME_BUCKET_NAME, bucketName); err != nil {
					t.Fatalf("Failed to set file RESUME_TOKEN_VOLUME_BUCKET_NAME environment variables.")
				}
				if err := setEnv(constant.RESUME_TOKEN_SAVE_INTERVAL_SEC, intavalSec); err != nil {
					t.Fatalf("Failed to set file RESUME_TOKEN_SAVE_INTERVAL_SEC environment variables.")
				}
				i, err := New(ctx, l)
				if err != nil {
					t.Fatalf("Testing Error, ErrorMessage: %v", err)
				}

				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				cli := mocks.NewMockStorageClient(ctrl)
				resumeToken, exists := i.(*resumeTokenImpl)
				if !exists {
					t.Fatalf("Testing Error. convert interaface to struct failed.")
				}
				resumeToken.client = cli

				ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
				defer cancel()

				cli.EXPECT().
					PutObject(ctx, fmt.Sprintf("%s/%s", resumeToken.volumePath, resumeToken.tokenFileName), rt).
					Return(nil).AnyTimes()

				env := resumeToken.Env()
				if env == "" {
					t.Fatalf("Failed to get environment variables for resume tokens settings.")
				}
				fmt.Printf("env %s\n", env)

				if err := resumeToken.SaveResumeToken(ctx, rt); err != nil {
					t.Fatalf("Testing Error, ErrorMessage: %v", err)
				}
				if err := resumeToken.SaveResumeToken(ctx, rt); err != nil {
					t.Fatalf("Testing Error, ErrorMessage: %v", err)
				}

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
