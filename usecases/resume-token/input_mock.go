//go:build test
// +build test

package resume_token

func (_ *mockResumeTokenClientImpl) fetchPersistentVolumeDir() (string, error) {
	return "", nil
}
