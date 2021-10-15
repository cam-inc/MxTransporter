package resume_token

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"mxtransporter/config"
	"mxtransporter/pkg/errors"
	"os"
	"time"
)

func SaveResumeToken(rt primitive.M) error {
	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		return errors.InternalServerError.Wrap("Failed to load location time.", err)
	}

	nowTime := time.Now().In(jst)
	nowYear := nowTime.Format("2006")
	nowMonth := nowTime.Format("01")
	nowDay := nowTime.Format("02")

	pv, err := config.PersistentVolume()
	if err != nil{
		return err
	}

	fileName := nowTime.Format("2006-01-02")
	filePath := pv + nowYear + "/" + nowMonth + "/" + nowDay + "/"
	file := filePath + fileName + ".dat"

	rtValue := rt["_data"].(string)

	if dirStat, err := os.Stat(filePath); os.IsNotExist(err) || dirStat.IsDir() {
		os.MkdirAll(filePath, 0777)
	}

	fp, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE, 0664)

	if err != nil {
		return errors.InternalServerErrorResumeToken.Wrap("Failed to open file saved resume token.", err)
	}
	defer fp.Close()

	_, err = fp.WriteString(rtValue)
	if err != nil {
		return errors.InternalServerErrorResumeToken.Wrap("Failed to write resume token in file.", err)
	}

	fmt.Println("Success to save a resume token in PVC")

	return nil
}
