package resume_token

import (
	"mxtransporter/config"
	"mxtransporter/pkg/errors"
	"mxtransporter/pkg/logger"
	"os"
	"time"
)

type (
	resumeTokenClient interface {
		fetchPersistentVolumeDir() (string, error)
	}

	ResumeTokenImpl struct {
		ResumeToken resumeTokenClient
		Log logger.Logger
	}

	ResumeTokenClientImpl struct{}

	mockResumeTokenClientImpl struct{}
)

func (_ *ResumeTokenClientImpl) fetchPersistentVolumeDir() (string, error) {
	pv, err := config.FetchPersistentVolumeDir()
	if err != nil {
		return "", err
	}
	return pv, nil
}

func (r *ResumeTokenImpl) SaveResumeToken(rt string) error {
	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		return errors.InternalServerError.Wrap("Failed to load location time.", err)
	}

	pv, err := r.ResumeToken.fetchPersistentVolumeDir()
	if err != nil {
		return err
	}

	nowTime := time.Now().In(jst)
	filePath := pv + nowTime.Format("2006/01/02/")
	file := filePath + nowTime.Format("2006-01-02.dat")

	if dirStat, err := os.Stat(filePath); os.IsNotExist(err) || dirStat.IsDir() {
		os.MkdirAll(filePath, 0777)
	}

	fp, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE, 0664)

	if err != nil {
		return errors.InternalServerError.Wrap("Failed to open file saved resume token.", err)
	}
	defer fp.Close()

	_, err = fp.WriteString(rt)
	if err != nil {
		return errors.InternalServerError.Wrap("Failed to write resume token in file.", err)
	}

	r.Log.ZLogger.Info("Success to save a resume token in PVC")

	return nil
}
