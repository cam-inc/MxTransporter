//go:build test
// +build test

package resume_token

import (
	"time"
)

type mockResumeTokenClientImpl struct{}

func (_ *mockResumeTokenClientImpl) fetchPersistentVolumeDir() (string, error) {
	return "", nil
}

func (_ *mockResumeTokenClientImpl) fetchNowTime() (time.Time, error) {
	// default time zone
	nowTime := time.Now()
	return nowTime, nil
}
