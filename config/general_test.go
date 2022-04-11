//go:build test
// +build test

package config

import (
	"github.com/cam-inc/mxtransporter/config/constant"
	"os"
	"reflect"
	"testing"
)

func Test_FetchExportDestination(t *testing.T) {
	tests := []struct {
		name   string
		runner func(t *testing.T)
	}{
		{
			name: "Does not succeed when the environment variable EXPORT_DESTINATION is not set.",
			runner: func(t *testing.T) {
				_, err := FetchExportDestination()
				if err == nil {
					t.Fatalf("Because the environment variable EXPORT_DESTINATION is not set, error must be returned from the target function.")
				}
			},
		},
		{
			name: "Check to call the set environment variable EXPORT_DESTINATION.",
			runner: func(t *testing.T) {
				expDst := "bigquery"
				if err := os.Setenv("EXPORT_DESTINATION", expDst); err != nil {
					t.Fatalf("Failed to set file EXPORT_DESTINATION environment variables.")
				}

				r, err := FetchExportDestination()
				if e, a := r, expDst; !reflect.DeepEqual(e, a) {
					t.Fatal("Environment variable EXPORT_DESTINATION is not acquired correctly.")
				}
				if err != nil {
					t.Fatal("Failed to fetch Environment variable EXPORT_DESTINATION.")
				}
			},
		},
	}

	for _, v := range tests {
		t.Run(v.name, v.runner)
	}
}

func Test_FetchGcpProject(t *testing.T) {
	tests := []struct {
		name   string
		runner func(t *testing.T)
	}{
		{
			name: "Does not succeed when the environment variable PROJECT_NAME_TO_EXPORT_CHANGE_STREAMS is not set.",
			runner: func(t *testing.T) {
				_, err := FetchGcpProject()
				if err == nil {
					t.Fatalf("Because the environment variable PROJECT_NAME_TO_EXPORT_CHANGE_STREAMS is not set, error must be returned from the target function.")
				}
			},
		},
		{
			name: "Check to call the set environment variable PROJECT_NAME_TO_EXPORT_CHANGE_STREAMS.",
			runner: func(t *testing.T) {
				projectID := "test-project"
				if err := os.Setenv("PROJECT_NAME_TO_EXPORT_CHANGE_STREAMS", projectID); err != nil {
					t.Fatalf("Failed to set file PROJECT_NAME_TO_EXPORT_CHANGE_STREAMS environment variables.")
				}

				r, err := FetchGcpProject()
				if e, a := r, projectID; !reflect.DeepEqual(e, a) {
					t.Fatal("Environment variable PROJECT_NAME_TO_EXPORT_CHANGE_STREAMS is not acquired correctly.")
				}
				if err != nil {
					t.Fatal("Failed to fetch Environment variable PROJECT_NAME_TO_EXPORT_CHANGE_STREAMS.")
				}
			},
		},
	}

	for _, v := range tests {
		t.Run(v.name, v.runner)
	}
}

func Test_FetchTimeZone(t *testing.T) {
	tests := []struct {
		name   string
		runner func(t *testing.T)
	}{
		{
			name: "Does not succeed when the environment variable TIME_ZONE is not set.",
			runner: func(t *testing.T) {
				_, err := FetchTimeZone()
				if err == nil {
					t.Fatalf("Because the environment variable TIME_ZONE is not set, error must be returned from the target function.")
				}
			},
		},
		{
			name: "Check to call the set environment variable TIME_ZONE.",
			runner: func(t *testing.T) {
				tz := "Asia/Tokyo"
				if err := os.Setenv("TIME_ZONE", tz); err != nil {
					t.Fatalf("Failed to set file TIME_ZONE environment variables.")
				}

				r, err := FetchTimeZone()
				if e, a := r, tz; !reflect.DeepEqual(e, a) {
					t.Fatal("Environment variable TIME_ZONE is not acquired correctly.")
				}
				if err != nil {
					t.Fatal("Failed to fetch Environment variable TIME_ZONE.")
				}
			},
		},
	}

	for _, v := range tests {
		t.Run(v.name, v.runner)
	}
}

func Test_LogConfig(t *testing.T) {
	t.Run("Check to call the set environment variable.", func(t *testing.T) {
		level := "1"
		format := "json"
		outputDir := "xxx"
		outputFile := "yyy"
		if err := os.Setenv("LOG_LEVEL", level); err != nil {
			t.Fatalf("Failed to set file LOG_LEVEL environment variables.")
		}
		if err := os.Setenv("LOG_FORMAT", format); err != nil {
			t.Fatalf("Failed to set file LOG_FORMAT environment variables.")
		}
		if err := os.Setenv("LOG_OUTPUT_DIRECTORY", outputDir); err != nil {
			t.Fatalf("Failed to set file LOG_OUTPUT_DIRECTORY environment variables.")
		}
		if err := os.Setenv("LOG_OUTPUT_FILE", outputFile); err != nil {
			t.Fatalf("Failed to set file LOG_OUTPUT_FILE environment variables.")
		}

		l := LogConfig()
		if e, a := l.Level, level; !reflect.DeepEqual(e, a) {
			t.Fatal("Environment variable LOG_LEVEL is not acquired correctly.")
		}
		if e, a := l.Format, format; !reflect.DeepEqual(e, a) {
			t.Fatal("Environment variable LOG_FORMAT is not acquired correctly.")
		}
		if e, a := l.OutputDirectory, outputDir; !reflect.DeepEqual(e, a) {
			t.Fatal("Environment variable LOG_OUTPUT_DIRECTORY is not acquired correctly.")
		}
		if e, a := l.OutputFile, outputFile; !reflect.DeepEqual(e, a) {
			t.Fatal("Environment variable LOG_OUTPUT_FILE is not acquired correctly.")
		}
	})
}

func Test_FetchResumeTokenFileName(t *testing.T) {
	t.Run("Check to call the set environment variable.", func(t *testing.T) {
		if err := os.Setenv("MONGODB_COLLECTION", "test"); err != nil {
			t.Fatalf("Failed to set file MONGODB_COLLECTION environment variables.")
		}
		col, err := FetchResumeTokenFileName()
		if err != nil {
			t.Fatalf("FetchResumeTokenFileName return error %v", err)
		}
		if e, a := col, "test.dat"; !reflect.DeepEqual(e, a) {
			t.Fatalf("Environment variable MONGODB_COLLECTION is not acquired correctly. %s", col)
		}

		if err := os.Setenv("RESUME_TOKEN_FILE_NAME", "resume"); err != nil {
			t.Fatalf("Failed to set file MONGO_COLLECTION environment variables.")
		}
		r, err := FetchResumeTokenFileName()
		if err != nil {
			t.Fatalf("FetchResumeTokenFileName return error %v", err)
		}
		if e, a := r, "resume"; !reflect.DeepEqual(e, a) {
			t.Fatalf("Environment variable LOG_OUTPUT_FILE is not acquired correctly. %s", r)
		}

		os.Unsetenv("RESUME_TOKEN_FILE_NAME")
		os.Unsetenv("MONGODB_COLLECTION")
		if _, err := FetchResumeTokenFileName(); err == nil {
			t.Fatal("FetchResumeTokenFileName no error.")
		}
	})
}

func Test_FileExportConfig(t *testing.T) {
	t.Run("Check to call the set environment variable.", func(t *testing.T) {
		if err := os.Setenv(constant.FILE_EXPORTER_CHANGE_STREAM_KEY, "changeStream"); err != nil {
			t.Fatalf("Failed to set file FILE_EXPORTER_CHANGE_STREAM_KEY environment variables.")
		}
		if err := os.Setenv(constant.FILE_EXPORTER_LOG_TYPE_KEY, "changeStream"); err != nil {
			t.Fatalf("Failed to set file FILE_EXPORTER_LOG_TYPE_KEY environment variables.")
		}
		if err := os.Setenv(constant.FILE_EXPORTER_TIME_KEY, "changeStream"); err != nil {
			t.Fatalf("Failed to set file FILE_EXPORTER_TIME_KEY environment variables.")
		}
		if err := os.Setenv(constant.FILE_EXPORTER_NAME_KEY, "changeStream"); err != nil {
			t.Fatalf("Failed to set file FILE_EXPORTER_NAME_KEY environment variables.")
		}
		FileExportConfig()
		os.Unsetenv(constant.FILE_EXPORTER_CHANGE_STREAM_KEY)
		os.Unsetenv(constant.FILE_EXPORTER_LOG_TYPE_KEY)
		os.Unsetenv(constant.FILE_EXPORTER_TIME_KEY)
		os.Unsetenv(constant.FILE_EXPORTER_NAME_KEY)
	})
}
