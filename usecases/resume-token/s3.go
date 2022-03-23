package resume_token

import (
	"bytes"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"go.uber.org/zap"
	"io"
	"path"
)

type (
	storageClient interface {
		GetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error)
		PutObject(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error)
	}

	rtS3 struct {
		log           *zap.SugaredLogger
		volumePath    string
		tokenFileName string
		client        storageClient
		bucket        string
	}
)

func (r *rtS3) ReadResumeToken(ctx context.Context) string {

	tmp := fmt.Sprintf("%s/%s", r.volumePath, r.tokenFileName)
	filePath := path.Clean(tmp)

	input := &s3.GetObjectInput{
		Bucket: aws.String(r.bucket),
		Key:    aws.String(filePath),
	}
	output, err := r.client.GetObject(ctx, input)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer output.Body.Close()
	o, err := io.ReadAll(output.Body)
	if err != nil {
		return ""
	}
	return string(o)
}

func (r *rtS3) SaveResumeToken(ctx context.Context, rt string) error {

	tmp := fmt.Sprintf("%s/%s", r.volumePath, r.tokenFileName)
	filePath := path.Clean(tmp)

	rtBuf := bytes.NewBuffer([]byte(rt))

	input := &s3.PutObjectInput{
		Bucket: aws.String(r.bucket),
		Key:    aws.String(filePath),
		Body:   rtBuf,
	}

	_, err := r.client.PutObject(ctx, input)
	if err != nil {
		return err
	}

	return nil
}

func (r *rtS3) Env() string {
	return fmt.Sprintf(`{"type": "s3", "volume_path":"%s", "file_name":"%s", "bucket":"%s"}`, r.volumePath, r.tokenFileName, r.bucket)
}
