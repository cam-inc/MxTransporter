//go:build test
// +build test

package resume_token

import (
	"fmt"
	"go.uber.org/zap"
	"mxtransporter/config"
	"mxtransporter/pkg/logger"
	"os"
	"reflect"
	"testing"
	"time"
)

func Test_SaveResumeToken(t *testing.T) {
	var l *zap.SugaredLogger
	logConfig := config.LogConfig()
	l = logger.New(logConfig)

	rt := "00000"

	tests := []struct {
		name   string
		runner func(t *testing.T)
	}{
		{
			name: "Failed to fetch PERSISTENT_VOLUME_DIR value.",
			runner: func(t *testing.T) {
				if err := os.Setenv("TIME_ZONE", "Asia/Tokyo"); err != nil {
					t.Fatalf("Failed to set file TIME_ZONE environment variables.")
				}

				resumeTokenImpl := ResumeTokenImpl{l}
				if err := resumeTokenImpl.SaveResumeToken(rt); err == nil {
					t.Fatalf("Not behaving as intended.")
				}

				// Reset the set environment variables for the next test.
				if err := os.Unsetenv("TIME_ZONE"); err != nil {
					t.Fatalf("Failed to unset file TIME_ZONE environment variables.")
				}
			},
		},
		{
			name: "Failed to fetch TIME_ZONE value.",
			runner: func(t *testing.T) {

				if err := os.Setenv("PERSISTENT_VOLUME_DIR", ""); err != nil {
					t.Fatalf("Failed to set file PERSISTENT_VOLUME_DIR environment variables.")
				}

				resumeTokenImpl := ResumeTokenImpl{l}
				if err := resumeTokenImpl.SaveResumeToken(rt); err == nil {
					t.Fatalf("Not behaving as intended.")
				}

				// Reset the set environment variables for the next test.
				if err := os.Unsetenv("PERSISTENT_VOLUME_DIR"); err != nil {
					t.Fatalf("Failed to unset file PERSISTENT_VOLUME_DIR environment variables.")
				}
			},
		},
		{
			name: "Pass to save resume token in file.",
			runner: func(t *testing.T) {
				if err := os.Setenv("PERSISTENT_VOLUME_DIR", ""); err != nil {
					t.Fatalf("Failed to set file PERSISTENT_VOLUME_DIR environment variables.")
				}

				if err := os.Setenv("TIME_ZONE", "Asia/Tokyo"); err != nil {
					t.Fatalf("Failed to set file TIME_ZONE environment variables.")
				}

				tl, err := time.LoadLocation("Asia/Tokyo")
				if err != nil {
					t.Fatalf("Failed to fetch time load location.")
				}

				nowTime := time.Now().In(tl)
				file := nowTime.Format("2006/01/02/2006-01-02.dat")

				resumeTokenImpl := ResumeTokenImpl{l}
				if err := resumeTokenImpl.SaveResumeToken(rt); err != nil {
					t.Fatalf("Testing Error, ErrorMessage: %v", err)
				}

				rtByte, err := os.ReadFile(file)
				if err != nil {
					t.Fatal("Failed to read file saved test resume token in.")
				}

				if e, a := rt, string(rtByte); !reflect.DeepEqual(e, a) {
					t.Errorf("expect %v, got %v", e, a)
				}
			},
		},
	}

	for _, v := range tests {
		t.Run(v.name, v.runner)
	}
}

func TestMain(m *testing.M) {
	status := m.Run()

	nowTime := time.Now()

	err := os.RemoveAll(nowTime.Format("2006"))
	if err != nil {
		fmt.Println(fmt.Errorf("The unnecessary file could not be deleted. errMessage: %s", err))
	}

	os.Exit(status)
}
