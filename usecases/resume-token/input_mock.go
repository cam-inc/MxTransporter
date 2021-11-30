//go:build test
// +build test

package resume_token

type mockResumeTokenClientImpl struct{}

func (_ *mockResumeTokenClientImpl) fetchPersistentVolumeDir() (string, error) {
	return "", nil
}
