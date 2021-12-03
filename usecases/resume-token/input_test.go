//go:build test
// +build test

package resume_token

import (
	"fmt"
	"go.uber.org/zap"
	"mxtransporter/config"
	"mxtransporter/pkg/errors"
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

	tl, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		t.Fatalf("Failed to fetch time load location.")
	}

	nowTime := time.Now().In(tl)
	file := nowTime.Format("2006/01/02/2006-01-02.dat")

	if err := os.Setenv("PERSISTENT_VOLUME_DIR", ""); err != nil {
		t.Fatalf("Failed to set file PERSISTENT_VOLUME_DIR environment variables.")
	}

	if err := os.Setenv("TIME_ZONE", "Asia/Tokyo"); err != nil {
		t.Fatalf("Failed to set file TIME_ZONE environment variables.")
	}

	rt := "00000"

	t.Run("Test if the resume token is stored in the correct location.", func(t *testing.T) {
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
	})
}

func TestMain(m *testing.M) {
	status := m.Run()

	nowTime := time.Now()

	err := os.RemoveAll(nowTime.Format("2006"))
	if err != nil {
		fmt.Println(errors.InternalServerError.Wrap("The unnecessary file could not be deleted.", err))
	}

	os.Exit(status)
}
