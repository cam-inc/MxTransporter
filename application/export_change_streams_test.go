//go:build test
// +build test

package application

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
	"mxtransporter/config"
	interfaceForBigquery "mxtransporter/interfaces/bigquery"
	interfaceForKinesisStream "mxtransporter/interfaces/kinesis-stream"
	interfaceForPubsub "mxtransporter/interfaces/pubsub"
	"mxtransporter/pkg/errors"
	"mxtransporter/pkg/logger"
	interfaceForResumeToken "mxtransporter/usecases/resume-token"
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
	ctx := context.TODO()

	var l *zap.SugaredLogger

	logCfg := config.LogConfig()
	l = logger.New(logCfg)

	if err := os.Setenv("TIME_ZONE", "Asia/Tokyo"); err != nil {
		t.Fatalf("Failed to set file TIME_ZONE environment variables.")
	}

	if err := os.Setenv("PERSISTENT_VOLUME_DIR", ""); err != nil {
		t.Fatalf("Failed to set file PERSISTENT_VOLUME_DIR environment variables.")
	}

	if err := os.Setenv("PROJECT_NAME_TO_EXPORT_CHANGE_STREAMS", ""); err != nil {
		t.Fatalf("Failed to set file GCP_PROJECT environment variables.")
	}

	tests := []struct {
		name   string
		runner func(t *testing.T)
	}{
		{
			name: "Pass to read resume token.",
			runner: func(t *testing.T) {
				if err := saveResumeToken("", "00000"); err != nil {
					t.Fatalf("Testing Error, ErrorMessage: %v", err)
				}
				if err := os.Setenv("EXPORT_DESTINATION", "bigquery"); err != nil {
					t.Fatalf("Failed to set file EXPORT_DESTINATION environment variables.")
				}
				mockWatcherClient := &mockChangeStremsWatcherClientImpl{nil, ChangeStreamsExporterImpl{}, "00000", false, "", "", ""}
				watcher := ChangeStremsWatcherImpl{mockWatcherClient, l}
				if err := watcher.WatchChangeStreams(ctx); err != nil {
					t.Fatalf("Testing Error, ErrorMessage: %v", err)
				}
				if mockWatcherClient.resumeAfterExistence == false {
					t.Fatalf("Testing Error, ErrorMessage: failed to get resumeToken.")
				}
			},
		},
		{
			name: "Pass not to read resume token.",
			runner: func(t *testing.T) {
				if err := os.Setenv("EXPORT_DESTINATION", "bigquery"); err != nil {
					t.Fatalf("Failed to set file EXPORT_DESTINATION environment variables.")
				}
				mockWatcherClient := &mockChangeStremsWatcherClientImpl{nil, ChangeStreamsExporterImpl{}, "", true, "", "", ""}
				watcher := ChangeStremsWatcherImpl{mockWatcherClient, l}
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
				if err := os.Setenv("EXPORT_DESTINATION", "bigquery"); err != nil {
					t.Fatalf("Failed to set file EXPORT_DESTINATION environment variables.")
				}
				mockWatcherClient := &mockChangeStremsWatcherClientImpl{nil, ChangeStreamsExporterImpl{}, "", true, "", "", ""}
				watcher := ChangeStremsWatcherImpl{mockWatcherClient, l}
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
				mockWatcherClient := &mockChangeStremsWatcherClientImpl{nil, ChangeStreamsExporterImpl{}, "", true, "", "", ""}
				watcher := ChangeStremsWatcherImpl{mockWatcherClient, l}
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
				if err := os.Setenv("EXPORT_DESTINATION", "kinesisStream"); err != nil {
					t.Fatalf("Failed to set file EXPORT_DESTINATION environment variables.")
				}
				mockWatcherClient := &mockChangeStremsWatcherClientImpl{nil, ChangeStreamsExporterImpl{}, "", true, "", "", ""}
				watcher := ChangeStremsWatcherImpl{mockWatcherClient, l}
				if err := watcher.WatchChangeStreams(ctx); err != nil {
					t.Fatalf("Testing Error, ErrorMessage: %v", err)
				}
				if mockWatcherClient.kinesisStreamPassCheck != "OK" {
					t.Fatalf("Testing Error, ErrorMessage: failed to get resumeToken.")
				}
			},
		},
		{
			name: "Pass to get bigquery, pubsub, kinesis stream client.",
			runner: func(t *testing.T) {
				if err := os.Setenv("EXPORT_DESTINATION", "bigquery,pubsub,kinesisStream"); err != nil {
					t.Fatalf("Failed to set file EXPORT_DESTINATION environment variables.")
				}
				mockWatcherClient := &mockChangeStremsWatcherClientImpl{nil, ChangeStreamsExporterImpl{}, "", true, "", "", ""}
				watcher := ChangeStremsWatcherImpl{mockWatcherClient, l}
				if err := watcher.WatchChangeStreams(ctx); err != nil {
					t.Fatalf("Testing Error, ErrorMessage: %v", err)
				}
				if mockWatcherClient.bqPassCheck != "OK" &&
					mockWatcherClient.pubsubPassCheck != "OK" &&
					mockWatcherClient.kinesisStreamPassCheck != "OK" {
					t.Fatalf("Testing Error, ErrorMessage: failed to get resumeToken.")
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
	ctx := context.TODO()

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

	mockExporterClient := &mockChangeStreamsExporterClientImpl{csMap, interfaceForBigquery.BigqueryImpl{}, interfaceForPubsub.PubsubImpl{}, interfaceForKinesisStream.KinesisStreamImpl{}, interfaceForResumeToken.ResumeTokenImpl{}, "", "", "", true}

	tests := []struct {
		name   string
		runner func(t *testing.T)
	}{
		{
			name: "Pass to export to bigquery.",
			runner: func(t *testing.T) {
				if err := os.Setenv("EXPORT_DESTINATION", "bigquery"); err != nil {
					t.Fatalf("Failed to set file EXPORT_DESTINATION environment variables.")
				}
				exporter := ChangeStreamsExporterImpl{mockExporterClient, l}
				if err := exporter.exportChangeStreams(ctx); err != nil {
					t.Fatalf("Testing Error, ErrorMessage: %v", err)
				}
				if mockExporterClient.bqPassCheck != "OK" {
					t.Fatalf("Testing Error, ErrorMessage: not going through export to bigquery.")
				}
			},
		},
		{
			name: "Pass to export to pubsub.",
			runner: func(t *testing.T) {
				if err := os.Setenv("EXPORT_DESTINATION", "pubsub"); err != nil {
					t.Fatalf("Failed to set file EXPORT_DESTINATION environment variables.")
				}
				exporter := ChangeStreamsExporterImpl{mockExporterClient, l}
				if err := exporter.exportChangeStreams(ctx); err != nil {
					t.Fatalf("Testing Error, ErrorMessage: %v", err)
				}
				if mockExporterClient.pubsubPassCheck != "OK" {
					t.Fatalf("Testing Error, ErrorMessage: not going through export to pubsub.")
				}
			},
		},
		{
			name: "Pass to export to kinesis stream.",
			runner: func(t *testing.T) {
				if err := os.Setenv("EXPORT_DESTINATION", "kinesisStream"); err != nil {
					t.Fatalf("Failed to set file EXPORT_DESTINATION environment variables.")
				}
				exporter := ChangeStreamsExporterImpl{mockExporterClient, l}
				if err := exporter.exportChangeStreams(ctx); err != nil {
					t.Fatalf("Testing Error, ErrorMessage: %v", err)
				}
				if mockExporterClient.kinesisStreamPassCheck != "OK" {
					t.Fatalf("Testing Error, ErrorMessage: not going through export to kinesis stream.")
				}
			},
		},
		{
			name: "Pass to export to kinesis stream.",
			runner: func(t *testing.T) {
				if err := os.Setenv("EXPORT_DESTINATION", "bigquery,pubsub,kinesisStream"); err != nil {
					t.Fatalf("Failed to set file EXPORT_DESTINATION environment variables.")
				}
				exporter := ChangeStreamsExporterImpl{mockExporterClient, l}
				if err := exporter.exportChangeStreams(ctx); err != nil {
					t.Fatalf("Testing Error, ErrorMessage: %v", err)
				}
				if mockExporterClient.bqPassCheck != "OK" &&
					mockExporterClient.pubsubPassCheck != "OK" &&
					mockExporterClient.kinesisStreamPassCheck != "OK" {
					t.Fatalf("Testing Error, ErrorMessage: not going through export to bigquery or pubsub or kinesis stream.")
				}
			},
		},
	}

	for _, v := range tests {
		t.Run(v.name, v.runner)
		mockExporterClient.csCursorFlag = true
	}
}
