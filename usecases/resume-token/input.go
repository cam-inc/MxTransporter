package resume_token

import (
	"go.uber.org/zap"
	"mxtransporter/config"
	"mxtransporter/pkg/common"
	"mxtransporter/pkg/errors"
	"os"
)

type ResumeTokenImpl struct {
	Log *zap.SugaredLogger
}

func (r *ResumeTokenImpl) SaveResumeToken(rt string) error {
	pv, err := config.FetchPersistentVolumeDir()
	if err != nil {
		return err
	}

	nowTime, err := common.FetchNowTime()
	if err != nil {
		return err
	}

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

	r.Log.Info("Success to save a resume token in PVC")

	return nil
}
