package resume_token

import (
	"context"
	"fmt"
	"github.com/cam-inc/mxtransporter/pkg/errors"
	"go.uber.org/zap"
	"os"
	"path"
)

type rtFile struct {
	log           *zap.SugaredLogger
	volumePath    string
	tokenFileName string
}

func (r *rtFile) ReadResumeToken(ctx context.Context) string {
	tmp := fmt.Sprintf("%s/%s", r.volumePath, r.tokenFileName)
	filePath := path.Clean(tmp)

	rtByte, err := os.ReadFile(filePath)
	if err != nil {
		return ""
	}

	return string(rtByte)

}

func (r *rtFile) SaveResumeToken(ctx context.Context, rt string) error {

	tmp := fmt.Sprintf("%s/%s", r.volumePath, r.tokenFileName)
	filePath := path.Clean(tmp)

	// os.IsNotExist(err) || dirStat.IsDir() { TODO:// なぜこうなってるのか質問する
	if _, err := os.Stat(r.volumePath); os.IsNotExist(err) {
		os.MkdirAll(r.volumePath, 0777)
	}

	fp, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0664)

	if err != nil {
		return errors.InternalServerError.Wrap("Failed to open file saved resume token.", err)
	}
	defer fp.Close()

	_, err = fp.WriteString(rt)
	if err != nil {
		return errors.InternalServerError.Wrap("Failed to write resume token in file.", err)
	}

	r.log.Info("Success to save a resume token in PVC")

	return nil
}

func (r *rtFile) Env() string {
	return fmt.Sprintf(`{"type": "file", "volume_path":"%s", "file_name":"%s"}`, r.volumePath, r.tokenFileName)
}
