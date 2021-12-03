//go:build test
// +build test

package logger

import (
	"mxtransporter/pkg/errors"
	"os"
	"testing"
)

var logCfg Log

func init() {
	logCfg.OutputDirectory = "test/"
	logCfg.OutputFile = "test.log"
}

func deleteFileSavedLog() error {
	err := os.RemoveAll(logCfg.OutputDirectory)
	if err != nil {
		return errors.InternalServerError.Wrap("The unnecessary file could not be deleted.", err)
	}
	return nil
}

func Test_New(t *testing.T) {
	tests := []struct {
		name   string
		runner func(t *testing.T)
	}{
		{
			name: "Check that the log saved to a file correctly. Info level logs should be output with default setting.",
			runner: func(t *testing.T) {
				l := New(logCfg)
				l.Info("test log")

				if _, err := os.Stat(logCfg.OutputDirectory); err != nil {
					t.Fatal("Failed to make directory.")
				}

				lf:= logCfg.OutputDirectory + logCfg.OutputFile
				if _, err := os.Stat(lf); err != nil {
					t.Fatal("Failed to make file.")
				}

				if lb , _ := os.ReadFile(lf); len(lb) == 0 {
					t.Fatal("Failed to save info log to log file.")
				}
			},
		},
		{
			name: "Check that the log saved to a file correctly.",
			runner: func(t *testing.T) {
				logCfg.Level = "0"
				l := New(logCfg)
				l.Info("test log")

				if _, err := os.Stat(logCfg.OutputDirectory); err != nil {
					t.Fatal("Failed to make directory.")
				}

				lf:= logCfg.OutputDirectory + logCfg.OutputFile
				if _, err := os.Stat(lf); err != nil {
					t.Fatal("Failed to make file.")
				}

				if lb , _ := os.ReadFile(lf); len(lb) == 0 {
					t.Fatal("Failed to save info log to log file.")
				}
			},
		},
		{
			name: "Check that the log output correctly.",
			runner: func(t *testing.T) {
				logCfg.Level = "1"
				l := New(logCfg)
				l.Info("test log")

				if _, err := os.Stat(logCfg.OutputDirectory); err != nil {
					t.Fatal("Failed to make directory.")
				}

				lf:= logCfg.OutputDirectory + logCfg.OutputFile
				if _, err := os.Stat(lf); err != nil {
					t.Fatal("Failed to make file.")
				}

				if lb , _ := os.ReadFile(lf); len(lb) != 0 {
					t.Fatal("Failed to output log. Info level logs should not be output in this test case.")
				}

				l.Error("test log")
				if lb , _ := os.ReadFile(lf); len(lb) == 0 {
					t.Fatal("Failed to save error log to log file.")
				}
			},
		},
	}
	for _, v := range tests {
		t.Run(v.name, v.runner)
		deleteFileSavedLog()
	}
}

