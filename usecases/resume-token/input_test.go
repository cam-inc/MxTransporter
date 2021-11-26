//go:build test
// +build test

package resume_token

import (
	"fmt"
	"mxtransporter/pkg/errors"
	"os"
	"reflect"
	"testing"
	"time"
)

func Test_SaveResumeToken(t *testing.T) {
	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		t.Fatal("Failed to load location time.")
	}

	nowTime := time.Now().In(jst)
	file := nowTime.Format("2006/01/02/2006-01-02.dat")

	rt := "00000"

	t.Run("Test if the resume token is stored in the correct location.", func(t *testing.T) {
		resumeTokenImpl := ResumeTokenImpl{&mockResumeTokenClientImpl{}}
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

	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		fmt.Println(errors.InternalServerError.Wrap("Failed to load location time.", err))
	}

	nowTime := time.Now().In(jst)

	err = os.RemoveAll(nowTime.Format("2006"))
	if err != nil {
		fmt.Println(errors.InternalServerError.Wrap("The unnecessary file could not be deleted.", err))
	}

	os.Exit(status)
}
