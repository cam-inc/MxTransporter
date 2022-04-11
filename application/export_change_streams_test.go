//go:build test
// +build test

package application

import (
	"context"
	mocks "github.com/cam-inc/mxtransporter/application/mock"
	"github.com/cam-inc/mxtransporter/config"
	interfaceForBigquery "github.com/cam-inc/mxtransporter/interfaces/bigquery"
	interfaceForKinesisStream "github.com/cam-inc/mxtransporter/interfaces/kinesis-stream"
	interfaceForPubsub "github.com/cam-inc/mxtransporter/interfaces/pubsub"
	"github.com/cam-inc/mxtransporter/pkg/errors"
	"github.com/cam-inc/mxtransporter/pkg/logger"
	"github.com/golang/mock/gomock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
	"os"
	"testing"
	"time"
)

func fetchNowTime(location string) (time.Time, error) {
	tl, err := time.LoadLocation(location)
	if err != nil {
		return time.Time{}, errors.InternalServerError.Wrap("Failed to fetch time load location.", err)
	}

	return time.Now().In(tl), nil
}

func saveResumeToken(pvDir string, rt string) error {
	nowTime, err := fetchNowTime("Asia/Tokyo")
	if err != nil {
		return errors.InternalServerError.Wrap("Failed to fetch now time.", err)
	}

	filePath := pvDir + nowTime.Format("2006/01/02/")
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
	return nil
}

func deleteFileSavedResumeToken() error {
	nowTime, err := fetchNowTime("Asia/Tokyo")
	if err != nil {
		return errors.InternalServerError.Wrap("Failed to fetch now time.", err)
	}

	err = os.RemoveAll(nowTime.Format("2006"))
	if err != nil {
		return errors.InternalServerError.Wrap("The unnecessary file could not be deleted.", err)
	}
	return nil
}

func Test_watchChangeStreams(t *testing.T) {
	ctx := context.Background()

	var l *zap.SugaredLogger

	logCfg := config.LogConfig()
	l = logger.New(logCfg)

	if err := os.Setenv("TIME_ZONE", "Asia/Tokyo"); err != nil {
		t.Fatalf("Failed to set file TIME_ZONE environment variables.")
	}

	if err := os.Setenv("PERSISTENT_VOLUME_DIR", ""); err != nil {
		t.Fatalf("Failed to set file PERSISTENT_VOLUME_DIR environment variables.")
	}

	if err := os.Setenv("EXPORT_DESTINATION", "bigquery"); err != nil {
		t.Fatalf("Failed to set file EXPORT_DESTINATION environment variables.")
	}

	if err := os.Setenv("PROJECT_NAME_TO_EXPORT_CHANGE_STREAMS", ""); err != nil {
		t.Fatalf("Failed to set file GCP_PROJECT environment variables.")
	}

	tests := []struct {
		name   string
		runner func(t *testing.T)
	}{
		{
			name: "Failed to fetch now time.",
			runner: func(t *testing.T) {
				// Unset environment variables to reproduce the condition.
				if err := os.Unsetenv("TIME_ZONE"); err != nil {
					t.Fatalf("Failed to unset file TIME_ZONE environment variables.")
				}
				mockWatcherClient := &mockChangeStremsWatcherClientImpl{
					mongoClient:            nil,
					csExporter:             ChangeStreamsExporterImpl{},
					resumeToken:            "",
					resumeAfterExistence:   false,
					bqPassCheck:            "",
					pubsubPassCheck:        "",
					kinesisStreamPassCheck: "",
					filePassCheck:          "",
				}
				watcher := ChangeStreamsWatcherImpl{
					Watcher: mockWatcherClient,
					Log:     l,
				}
				if err := watcher.WatchChangeStreams(ctx); err == nil {
					t.Fatalf("Not behaving as intended.")
				}

				// Undo environment variables
				if err := os.Setenv("TIME_ZONE", "Asia/Tokyo"); err != nil {
					t.Fatalf("Failed to set file TIME_ZONE environment variables.")
				}
			},
		},
		{
			name: "Failed to fetch persistent volume directory.",
			runner: func(t *testing.T) {
				// Unset environment variables to reproduce the condition.
				if err := os.Unsetenv("PERSISTENT_VOLUME_DIR"); err != nil {
					t.Fatalf("Failed to unset file PERSISTENT_VOLUME_DIR environment variables.")
				}

				mockWatcherClient := &mockChangeStremsWatcherClientImpl{
					mongoClient:            nil,
					csExporter:             ChangeStreamsExporterImpl{},
					resumeToken:            "",
					resumeAfterExistence:   false,
					bqPassCheck:            "",
					pubsubPassCheck:        "",
					kinesisStreamPassCheck: "",
					filePassCheck:          "",
				}
				watcher := ChangeStreamsWatcherImpl{
					Watcher: mockWatcherClient,
					Log:     l,
				}
				if err := watcher.WatchChangeStreams(ctx); err == nil {
					t.Fatalf("Not behaving as intended.")
				}

				// Undo environment variables
				if err := os.Setenv("PERSISTENT_VOLUME_DIR", ""); err != nil {
					t.Fatalf("Failed to set file PERSISTENT_VOLUME_DIR environment variables.")
				}
			},
		},
		{
			name: "Failed to fetch export destination.",
			runner: func(t *testing.T) {
				// Unset environment variables to reproduce the condition.
				if err := os.Unsetenv("EXPORT_DESTINATION"); err != nil {
					t.Fatalf("Failed to unset file EXPORT_DESTINATION environment variables.")
				}

				mockWatcherClient := &mockChangeStremsWatcherClientImpl{
					mongoClient:            nil,
					csExporter:             ChangeStreamsExporterImpl{},
					resumeToken:            "",
					resumeAfterExistence:   false,
					bqPassCheck:            "",
					pubsubPassCheck:        "",
					kinesisStreamPassCheck: "",
					filePassCheck:          "",
				}
				watcher := ChangeStreamsWatcherImpl{
					Watcher: mockWatcherClient,
					Log:     l,
				}
				if err := watcher.WatchChangeStreams(ctx); err == nil {
					t.Fatalf("Not behaving as intended.")
				}

				// Undo environment variables
				if err := os.Setenv("EXPORT_DESTINATION", "bigquery"); err != nil {
					t.Fatalf("Failed to set file EXPORT_DESTINATION environment variables.")
				}
			},
		},
		{
			name: "Failed to fetch gcp project id.",
			runner: func(t *testing.T) {
				// Unset environment variables to reproduce the condition.
				if err := os.Unsetenv("PROJECT_NAME_TO_EXPORT_CHANGE_STREAMS"); err != nil {
					t.Fatalf("Failed to unset file PROJECT_NAME_TO_EXPORT_CHANGE_STREAMS environment variables.")
				}

				mockWatcherClient := &mockChangeStremsWatcherClientImpl{
					mongoClient:            nil,
					csExporter:             ChangeStreamsExporterImpl{},
					resumeToken:            "",
					resumeAfterExistence:   false,
					bqPassCheck:            "",
					pubsubPassCheck:        "",
					kinesisStreamPassCheck: "",
					filePassCheck:          "",
				}
				watcher := ChangeStreamsWatcherImpl{
					Watcher: mockWatcherClient,
					Log:     l,
				}
				if err := watcher.WatchChangeStreams(ctx); err == nil {
					t.Fatalf("Not behaving as intended.")
				}

				// Undo environment variables
				if err := os.Setenv("PROJECT_NAME_TO_EXPORT_CHANGE_STREAMS", ""); err != nil {
					t.Fatalf("Failed to set file GCP_PROJECT environment variables.")
				}
			},
		},
		{
			name: "Pass to read resume token.",
			runner: func(t *testing.T) {
				token := "00000"

				if err := os.Setenv("MONGODB_COLLECTION", "test"); err != nil {
					t.Fatalf("Failed to set file MONGODB_COLLECTION environment variables.")
				}
				mockWatcherClient := &mockChangeStremsWatcherClientImpl{
					mongoClient:            nil,
					csExporter:             ChangeStreamsExporterImpl{},
					resumeToken:            "00000",
					resumeAfterExistence:   false,
					bqPassCheck:            "",
					pubsubPassCheck:        "",
					kinesisStreamPassCheck: "",
					filePassCheck:          "",
				}
				watcher := ChangeStreamsWatcherImpl{
					Watcher: mockWatcherClient,
					Log:     l,
				}

				ctrl := gomock.NewController(t)
				defer ctrl.Finish()
				resumeTokenImpl := mocks.NewMockResumeToken(ctrl)
				resumeTokenImpl.EXPECT().ReadResumeToken(ctx).Return(token).AnyTimes()
				watcher.setResumeTokenManager(resumeTokenImpl)

				if err := watcher.WatchChangeStreams(ctx); err != nil {
					t.Fatalf("Testing Error, ErrorMessage: %v", err)
				}
				if mockWatcherClient.resumeAfterExistence == false {
					t.Fatalf("Testing Error, ErrorMessage: failed to get resumeToken.")
				}
			},
		},
		{
			name: "Failed to read resume token.",
			runner: func(t *testing.T) {
				mockWatcherClient := &mockChangeStremsWatcherClientImpl{
					mongoClient:            nil,
					csExporter:             ChangeStreamsExporterImpl{},
					resumeToken:            "",
					resumeAfterExistence:   true,
					bqPassCheck:            "",
					pubsubPassCheck:        "",
					kinesisStreamPassCheck: "",
					filePassCheck:          "",
				}
				watcher := ChangeStreamsWatcherImpl{
					Watcher: mockWatcherClient,
					Log:     l,
				}
				if err := watcher.WatchChangeStreams(ctx); err != nil {
					t.Fatalf("Testing Error, ErrorMessage: %v", err)
				}
				if mockWatcherClient.resumeAfterExistence == true {
					t.Fatalf("Testing Error, ErrorMessage: failed to get resumeToken.")
				}
			},
		},
		{
			name: "Pass to get bigquery client.",
			runner: func(t *testing.T) {
				mockWatcherClient := &mockChangeStremsWatcherClientImpl{
					mongoClient:            nil,
					csExporter:             ChangeStreamsExporterImpl{},
					resumeToken:            "",
					resumeAfterExistence:   true,
					bqPassCheck:            "",
					pubsubPassCheck:        "",
					kinesisStreamPassCheck: "",
					filePassCheck:          "",
				}
				watcher := ChangeStreamsWatcherImpl{
					Watcher: mockWatcherClient,
					Log:     l,
				}
				if err := watcher.WatchChangeStreams(ctx); err != nil {
					t.Fatalf("Testing Error, ErrorMessage: %v", err)
				}
				if mockWatcherClient.bqPassCheck != "OK" {
					t.Fatalf("Testing Error, ErrorMessage: failed to get resumeToken.")
				}
			},
		},
		{
			name: "Pass to get pubsub client.",
			runner: func(t *testing.T) {
				if err := os.Setenv("EXPORT_DESTINATION", "pubsub"); err != nil {
					t.Fatalf("Failed to set file EXPORT_DESTINATION environment variables.")
				}
				mockWatcherClient := &mockChangeStremsWatcherClientImpl{
					mongoClient:            nil,
					csExporter:             ChangeStreamsExporterImpl{},
					resumeToken:            "",
					resumeAfterExistence:   true,
					bqPassCheck:            "",
					pubsubPassCheck:        "",
					kinesisStreamPassCheck: "",
					filePassCheck:          "",
				}
				watcher := ChangeStreamsWatcherImpl{
					Watcher: mockWatcherClient,
					Log:     l,
				}
				if err := watcher.WatchChangeStreams(ctx); err != nil {
					t.Fatalf("Testing Error, ErrorMessage: %v", err)
				}
				if mockWatcherClient.pubsubPassCheck != "OK" {
					t.Fatalf("Testing Error, ErrorMessage: failed to get resumeToken.")
				}
			},
		},
		{
			name: "Pass to get kinesis stream client.",
			runner: func(t *testing.T) {
				// Unset environment variables to reproduce the condition.
				if err := os.Unsetenv("PROJECT_NAME_TO_EXPORT_CHANGE_STREAMS"); err != nil {
					t.Fatalf("Failed to unset file PROJECT_NAME_TO_EXPORT_CHANGE_STREAMS environment variables.")
				}

				if err := os.Setenv("EXPORT_DESTINATION", "kinesisStream"); err != nil {
					t.Fatalf("Failed to set file EXPORT_DESTINATION environment variables.")
				}
				mockWatcherClient := &mockChangeStremsWatcherClientImpl{
					mongoClient:            nil,
					csExporter:             ChangeStreamsExporterImpl{},
					resumeToken:            "",
					resumeAfterExistence:   true,
					bqPassCheck:            "",
					pubsubPassCheck:        "",
					kinesisStreamPassCheck: "",
					filePassCheck:          "",
				}
				watcher := ChangeStreamsWatcherImpl{
					Watcher: mockWatcherClient,
					Log:     l,
				}
				if err := watcher.WatchChangeStreams(ctx); err != nil {
					t.Fatalf("Testing Error, ErrorMessage: %v", err)
				}
				if mockWatcherClient.kinesisStreamPassCheck != "OK" {
					t.Fatalf("Testing Error, ErrorMessage: failed to get resumeToken.")
				}

				// Undo environment variables
				if err := os.Setenv("PROJECT_NAME_TO_EXPORT_CHANGE_STREAMS", ""); err != nil {
					t.Fatalf("Failed to set file GCP_PROJECT environment variables.")
				}
			},
		},
		{
			name: "Pass to get file output client.",
			runner: func(t *testing.T) {
				// Unset environment variables to reproduce the condition.
				if err := os.Unsetenv("PROJECT_NAME_TO_EXPORT_CHANGE_STREAMS"); err != nil {
					t.Fatalf("Failed to unset file PROJECT_NAME_TO_EXPORT_CHANGE_STREAMS environment variables.")
				}

				if err := os.Setenv("EXPORT_DESTINATION", "file"); err != nil {
					t.Fatalf("Failed to set file EXPORT_DESTINATION environment variables.")
				}
				mockWatcherClient := &mockChangeStremsWatcherClientImpl{
					mongoClient:            nil,
					csExporter:             ChangeStreamsExporterImpl{},
					resumeToken:            "",
					resumeAfterExistence:   true,
					bqPassCheck:            "",
					pubsubPassCheck:        "",
					kinesisStreamPassCheck: "",
					filePassCheck:          "",
				}
				watcher := ChangeStreamsWatcherImpl{
					Watcher: mockWatcherClient,
					Log:     l,
				}
				if err := watcher.WatchChangeStreams(ctx); err != nil {
					t.Fatalf("Testing Error, ErrorMessage: %v", err)
				}
				if mockWatcherClient.filePassCheck != "OK" {
					t.Fatalf("Testing Error, ErrorMessage: failed to get resumeToken.")
				}

				// Undo environment variables
				if err := os.Setenv("PROJECT_NAME_TO_EXPORT_CHANGE_STREAMS", ""); err != nil {
					t.Fatalf("Failed to set file GCP_PROJECT environment variables.")
				}
			},
		},
		{
			name: "Pass to get bigquery, pubsub, kinesis stream, file client.",
			runner: func(t *testing.T) {
				if err := os.Setenv("EXPORT_DESTINATION", "bigquery,pubsub,kinesisStream,file"); err != nil {
					t.Fatalf("Failed to set file EXPORT_DESTINATION environment variables.")
				}
				mockWatcherClient := &mockChangeStremsWatcherClientImpl{
					mongoClient:            nil,
					csExporter:             ChangeStreamsExporterImpl{},
					resumeToken:            "",
					resumeAfterExistence:   true,
					bqPassCheck:            "",
					pubsubPassCheck:        "",
					kinesisStreamPassCheck: "",
					filePassCheck:          "",
				}
				watcher := ChangeStreamsWatcherImpl{
					Watcher: mockWatcherClient,
					Log:     l,
				}
				if err := watcher.WatchChangeStreams(ctx); err != nil {
					t.Fatalf("Testing Error, ErrorMessage: %v", err)
				}
				if mockWatcherClient.bqPassCheck != "OK" &&
					mockWatcherClient.pubsubPassCheck != "OK" &&
					mockWatcherClient.kinesisStreamPassCheck != "OK" &&
					mockWatcherClient.filePassCheck != "OK" {
					t.Fatalf("Testing Error, ErrorMessage: failed to get multi client.")
				}
			},
		},
		{
			name: "Export destination is wrong.",
			runner: func(t *testing.T) {
				if err := os.Setenv("EXPORT_DESTINATION", "xxx"); err != nil {
					t.Fatalf("Failed to set file EXPORT_DESTINATION environment variables.")
				}
				mockWatcherClient := &mockChangeStremsWatcherClientImpl{
					mongoClient:            nil,
					csExporter:             ChangeStreamsExporterImpl{},
					resumeToken:            "",
					resumeAfterExistence:   true,
					bqPassCheck:            "",
					pubsubPassCheck:        "",
					kinesisStreamPassCheck: "",
					filePassCheck:          "",
				}
				watcher := ChangeStreamsWatcherImpl{
					Watcher: mockWatcherClient,
					Log:     l,
				}
				if err := watcher.WatchChangeStreams(ctx); err == nil {
					t.Fatalf("Not behaving as intended.")
				}
			},
		},
	}

	for _, v := range tests {
		t.Run(v.name, v.runner)
		if err := deleteFileSavedResumeToken(); err != nil {
			t.Fatalf("Testing Error, ErrorMessage: %v", err)
		}
	}
}

func Test_exportChangeStreams(t *testing.T) {
	ctx := context.Background()

	var l *zap.SugaredLogger

	logCfg := config.LogConfig()
	l = logger.New(logCfg)

	csMap := primitive.M{
		"ns": primitive.M{
			"db":   "test db",
			"coll": "test coll",
		},
		"operationType": "insert",
		"clusterTime": primitive.Timestamp{
			T: 1638284400,
			I: 1,
		},
		"_id": primitive.M{
			"_data": "00000",
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	resumeTokenImpl := mocks.NewMockResumeToken(ctrl)

	mockExporterClient := &mockChangeStreamsExporterClientImpl{
		cs:                     csMap,
		bq:                     interfaceForBigquery.BigqueryImpl{},
		pubsub:                 interfaceForPubsub.PubsubImpl{},
		kinesisStream:          interfaceForKinesisStream.KinesisStreamImpl{},
		resumeToken:            resumeTokenImpl,
		bqPassCheck:            "",
		pubsubPassCheck:        "",
		kinesisStreamPassCheck: "",
		filePassCheck:          "",
		csCursorFlag:           true,
	}

	tests := []struct {
		name   string
		runner func(t *testing.T)
	}{
		{
			name: "Failed to fetch export destination.",
			runner: func(t *testing.T) {
				exporter := ChangeStreamsExporterImpl{mockExporterClient, l}
				if err := exporter.exportChangeStreams(ctx); err == nil {
					t.Fatalf("Not behaving as intended.")
				}
			},
		},
		{
			name: "Pass to export to bigquery.",
			runner: func(t *testing.T) {
				if err := os.Setenv("EXPORT_DESTINATION", "bigquery"); err != nil {
					t.Fatalf("Failed to set file EXPORT_DESTINATION environment variables.")
				}
				if err := os.Setenv("MONGODB_COLLECTION", "test"); err != nil {
					t.Fatalf("Failed to set file MONGODB_COLLECTION environment variables.")
				}
				exporter := ChangeStreamsExporterImpl{mockExporterClient, l}
				if err := exporter.exportChangeStreams(ctx); err != nil {
					t.Fatalf("Testing Error, ErrorMessage: %v", err)
				}
				if mockExporterClient.bqPassCheck != "OK" {
					t.Fatalf("Testing Error, ErrorMessage: not going through export to bigquery.")
				}
				mockExporterClient.bqPassCheck = ""
			},
		},
		{
			name: "Pass to export to pubsub.",
			runner: func(t *testing.T) {
				if err := os.Setenv("EXPORT_DESTINATION", "pubsub"); err != nil {
					t.Fatalf("Failed to set file EXPORT_DESTINATION environment variables.")
				}
				if err := os.Setenv("MONGODB_COLLECTION", "test"); err != nil {
					t.Fatalf("Failed to set file MONGODB_COLLECTION environment variables.")
				}
				exporter := ChangeStreamsExporterImpl{mockExporterClient, l}
				if err := exporter.exportChangeStreams(ctx); err != nil {
					t.Fatalf("Testing Error, ErrorMessage: %v", err)
				}
				if mockExporterClient.pubsubPassCheck != "OK" {
					t.Fatalf("Testing Error, ErrorMessage: not going through export to pubsub.")
				}
				mockExporterClient.pubsubPassCheck = ""
			},
		},
		{
			name: "Pass to export to kinesis stream.",
			runner: func(t *testing.T) {
				if err := os.Setenv("EXPORT_DESTINATION", "kinesisStream"); err != nil {
					t.Fatalf("Failed to set file EXPORT_DESTINATION environment variables.")
				}
				if err := os.Setenv("MONGODB_COLLECTION", "test"); err != nil {
					t.Fatalf("Failed to set file MONGODB_COLLECTION environment variables.")
				}
				exporter := ChangeStreamsExporterImpl{mockExporterClient, l}
				if err := exporter.exportChangeStreams(ctx); err != nil {
					t.Fatalf("Testing Error, ErrorMessage: %v", err)
				}
				if mockExporterClient.kinesisStreamPassCheck != "OK" {
					t.Fatalf("Testing Error, ErrorMessage: not going through export to kinesis stream.")
				}
				mockExporterClient.kinesisStreamPassCheck = ""
			},
		},
		{
			name: "Pass to export to file.",
			runner: func(t *testing.T) {
				if err := os.Setenv("EXPORT_DESTINATION", "file"); err != nil {
					t.Fatalf("Failed to set file EXPORT_DESTINATION environment variables.")
				}
				if err := os.Setenv("MONGODB_COLLECTION", "test"); err != nil {
					t.Fatalf("Failed to set file MONGODB_COLLECTION environment variables.")
				}
				exporter := ChangeStreamsExporterImpl{mockExporterClient, l}
				if err := exporter.exportChangeStreams(ctx); err != nil {
					t.Fatalf("Testing Error, ErrorMessage: %v", err)
				}
				if mockExporterClient.filePassCheck != "OK" {
					t.Fatalf("Testing Error, ErrorMessage: not going through export to local storage file.")
				}
				mockExporterClient.filePassCheck = ""
			},
		},
		{
			name: "Pass to export to bigquery, pubsub, kinesis stream, file.",
			runner: func(t *testing.T) {
				if err := os.Setenv("EXPORT_DESTINATION", "bigquery,pubsub,kinesisStream,file"); err != nil {
					t.Fatalf("Failed to set file EXPORT_DESTINATION environment variables.")
				}
				if err := os.Setenv("MONGODB_COLLECTION", "test"); err != nil {
					t.Fatalf("Failed to set file MONGODB_COLLECTION environment variables.")
				}
				exporter := ChangeStreamsExporterImpl{mockExporterClient, l}
				if err := exporter.exportChangeStreams(ctx); err != nil {
					t.Fatalf("Testing Error, ErrorMessage: %v", err)
				}
				if mockExporterClient.bqPassCheck != "OK" &&
					mockExporterClient.pubsubPassCheck != "OK" &&
					mockExporterClient.kinesisStreamPassCheck != "OK" &&
					mockExporterClient.filePassCheck != "OK" {
					t.Fatalf("Testing Error, ErrorMessage: not going through export to bigquery or pubsub or kinesis stream.")
				}
			},
		},
		{
			name: "Export destination is wrong.",
			runner: func(t *testing.T) {
				if err := os.Setenv("EXPORT_DESTINATION", "xxx"); err != nil {
					t.Fatalf("Failed to set file EXPORT_DESTINATION environment variables.")
				}
				if err := os.Setenv("MONGODB_COLLECTION", "test"); err != nil {
					t.Fatalf("Failed to set file MONGODB_COLLECTION environment variables.")
				}
				exporter := ChangeStreamsExporterImpl{mockExporterClient, l}
				if err := exporter.exportChangeStreams(ctx); err == nil {
					t.Fatalf("Not behaving as intended.")
				}
			},
		},
	}

	for _, v := range tests {
		t.Run(v.name, v.runner)
		mockExporterClient.csCursorFlag = true
	}
}
