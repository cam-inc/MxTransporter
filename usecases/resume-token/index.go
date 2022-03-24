package resume_token

import (
	"context"
	"fmt"
	"github.com/cam-inc/mxtransporter/config"
	rtConfig "github.com/cam-inc/mxtransporter/config/resume-token"
	"github.com/cam-inc/mxtransporter/pkg/client"
	"github.com/cam-inc/mxtransporter/pkg/client/storage"
	"github.com/cam-inc/mxtransporter/pkg/errors"
	"go.uber.org/zap"
	"path"
	"sync"
	"time"
)

type ResumeToken interface {
	ReadResumeToken(ctx context.Context) string
	SaveResumeToken(ctx context.Context, rt string) error
	Env() string
}

type resumeTokenImpl struct {
	Log             *zap.SugaredLogger
	client          storage.StorageClient
	volumeType      string
	volumePath      string
	tokenFileName   string
	saveIntervalSec int
	savedTimestamp  time.Time
	lock            sync.Locker
}

func (r *resumeTokenImpl) ReadResumeToken(ctx context.Context) string {
	tmp := fmt.Sprintf("%s/%s", r.volumePath, r.tokenFileName)
	filePath := path.Clean(tmp)
	o, err := r.client.GetObject(ctx, filePath)
	if err != nil {
		r.Log.Errorf("Failed ReadResumeToken key:%s, err:%v", filePath, err)
		return ""
	}
	return string(o)
}

func (r *resumeTokenImpl) SaveResumeToken(ctx context.Context, rt string) error {
	if !r.enableSave() {
		return nil
	}
	tmp := fmt.Sprintf("%s/%s", r.volumePath, r.tokenFileName)
	filePath := path.Clean(tmp)
	if err := r.client.PutObject(ctx, filePath, rt); err != nil {
		r.Log.Errorf("Failed SaveResumeToken key:%s, err:%v", filePath, err)
		return errors.InternalServerError.Wrap("Failed to SaveResumeToken", err)
	}
	r.setSavedTimestamp()
	return nil
}

func (r *resumeTokenImpl) setSavedTimestamp() {
	if r.saveIntervalSec == 0 {
		return
	}
	r.lock.Lock()
	defer r.lock.Unlock()
	r.savedTimestamp = time.Now()
}
func (r *resumeTokenImpl) enableSave() bool {
	if r.saveIntervalSec == 0 {
		return true
	}
	t := r.savedTimestamp.Add(time.Duration(r.saveIntervalSec) * time.Second)
	return t.Before(time.Now())
}

func (r resumeTokenImpl) Env() string {
	return fmt.Sprintf(`{"volumeType":"%s","volume_path":"%s", "file_name":"%s", "intaval":"%d"}`, r.volumeType, r.volumePath, r.tokenFileName, r.saveIntervalSec)
}

func New(ctx context.Context, log *zap.SugaredLogger) (ResumeToken, error) {

	cfg := rtConfig.ResumeTokenConfig()

	fileName, err := config.FetchResumeTokenFileName()
	if err != nil {
		return nil, err
	}

	cli, err := client.NewResumeTokenClient(ctx, cfg)
	if err != nil {
		return nil, err
	}

	mu := &sync.RWMutex{}
	return &resumeTokenImpl{
		Log:             log,
		volumeType:      cfg.VolumeType,
		volumePath:      cfg.Path,
		tokenFileName:   fileName,
		saveIntervalSec: cfg.SaveIntervalSec,
		lock:            mu.RLocker(),
		client:          cli,
	}, nil
}
