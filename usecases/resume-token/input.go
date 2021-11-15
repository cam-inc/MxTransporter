package resume_token

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"mxtransporter/config"
	"mxtransporter/pkg/errors"
	"os"
	"time"
)

type generalConfig struct {
	gerneralConfigIf config.GeneralConfigIf
}

func NewGeneralConfig(gerneralConfigIf config.GeneralConfigIf) *generalConfig {
	return &generalConfig{
		gerneralConfigIf: gerneralConfigIf,
	}
}

func (c *generalConfig) SaveResumeToken(rt primitive.M) error {
	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		return errors.InternalServerError.Wrap("Failed to load location time.", err)
	}

	nowTime := time.Now().In(jst)

	pv, err := c.gerneralConfigIf.FetchPersistentVolumeDir()
	if err != nil {
		return err
	}

	filePath := pv + nowTime.Format("2006/01/02/")
	file := filePath + nowTime.Format("2006-01-02.dat")

	rtValue := rt["_data"].(string)

	if dirStat, err := os.Stat(filePath); os.IsNotExist(err) || dirStat.IsDir() {
		os.MkdirAll(filePath, 0777)
	}

	fp, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE, 0664)

	if err != nil {
		return errors.InternalServerError.Wrap("Failed to open file saved resume token.", err)
	}
	defer fp.Close()

	_, err = fp.WriteString(rtValue)
	if err != nil {
		return errors.InternalServerError.Wrap("Failed to write resume token in file.", err)
	}

	fmt.Println("Success to save a resume token in PVC")

	return nil
}
