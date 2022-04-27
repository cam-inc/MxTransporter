//go:build test
// +build test

package logger

import (
	"bytes"
	"github.com/cam-inc/mxtransporter/pkg/errors"
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
			name: "Check that a log is saved to a file correctly. Info level logs should be output with default setting.",
			runner: func(t *testing.T) {
				l := New(logCfg)
				l.Info("test log")

				if _, err := os.Stat(logCfg.OutputDirectory); err != nil {
					t.Fatal("Failed to make directory.")
				}

				lf := logCfg.OutputDirectory + logCfg.OutputFile
				if _, err := os.Stat(lf); err != nil {
					t.Fatal("Failed to make file.")
				}

				if lb, _ := os.ReadFile(lf); len(lb) == 0 {
					t.Fatal("Failed to save info log to log file.")
				}
			},
		},
		{
			name: "Check that a log is saved to a file correctly.",
			runner: func(t *testing.T) {
				logCfg.Level = "0"
				l := New(logCfg)
				l.Info("test log")

				if _, err := os.Stat(logCfg.OutputDirectory); err != nil {
					t.Fatal("Failed to make directory.")
				}

				lf := logCfg.OutputDirectory + logCfg.OutputFile
				if _, err := os.Stat(lf); err != nil {
					t.Fatal("Failed to make file.")
				}

				if lb, _ := os.ReadFile(lf); len(lb) == 0 {
					t.Fatal("Failed to save info log to log file.")
				}
			},
		},
		{
			name: "Pass to output a json log to file correctly.",
			runner: func(t *testing.T) {
				logCfg.Level = "1"
				l := New(logCfg)
				l.Info("test log")

				if _, err := os.Stat(logCfg.OutputDirectory); err != nil {
					t.Fatal("Failed to make directory.")
				}

				lf := logCfg.OutputDirectory + logCfg.OutputFile
				if _, err := os.Stat(lf); err != nil {
					t.Fatal("Failed to make file.")
				}

				if lb, _ := os.ReadFile(lf); len(lb) != 0 {
					t.Fatal("Failed to output log. Info level logs should not be output in this test case.")
				}

				l.Error("test log")
				if lb, _ := os.ReadFile(lf); len(lb) == 0 {
					t.Fatal("Failed to save error log to log file.")
				}
			},
		},
		{
			name: "Pass to output a console log to file correctly.",
			runner: func(t *testing.T) {
				logCfg.Level = "1"
				logCfg.Format = "console"
				l := New(logCfg)
				l.Info("test log")

				if _, err := os.Stat(logCfg.OutputDirectory); err != nil {
					t.Fatal("Failed to make directory.")
				}

				lf := logCfg.OutputDirectory + logCfg.OutputFile
				if _, err := os.Stat(lf); err != nil {
					t.Fatal("Failed to make file.")
				}

				if lb, _ := os.ReadFile(lf); len(lb) != 0 {
					t.Fatal("Failed to output log. Info level logs should not be output in this test case.")
				}

				l.Error("test log")
				if lb, _ := os.ReadFile(lf); len(lb) == 0 {
					t.Fatal("Failed to save error log to log file.")
				}
			},
		},
		{
			name: "Pass to output a log to stdout.",
			runner: func(t *testing.T) {
				stdout := os.Stdout

				r, w, err := os.Pipe()
				if err != nil {
					t.Fatal("Failed to create pipe of w/r file.")
				}

				os.Stdout = w

				logCfgDefault := Log{}
				l := New(logCfgDefault)
				l.Info("test log")

				w.Close()

				os.Stdout = stdout

				var buf bytes.Buffer
				buf.ReadFrom(r)

				// Even if it is not written to the stdout, 16 characters will be passed when passing the value from the reader of os.Pipe () to the buffer.
				// so make the following conditions.
				if len(buf.String()) < 17 {
					t.Fatal("Failed to save a log into stdout.")
				}
			},
		},
	}
	for _, v := range tests {
		t.Run(v.name, v.runner)
		deleteFileSavedLog()
	}
}
