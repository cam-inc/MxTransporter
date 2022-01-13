package common

import (
	"os"
	"testing"
	"time"
)

func Test_Contains(t *testing.T) {
	tests := []struct {
		name   string
		runner func(t *testing.T)
	}{
		{
			name: "Check that the function works correctly.",
			runner: func(t *testing.T) {
				a := []string{"Test", "is", "most", "important", "program"}
				s := "Test"
				ok := Contains(a, s)
				if ok != true {
					t.Fatal("The function is not working properly.")
				}
			},
		},
		{
			name: "Check that the function fail.",
			runner: func(t *testing.T) {
				a := []string{"Test", "is", "most", "important", "program"}
				s := "xxx"
				ok := Contains(a, s)
				if ok != false {
					t.Fatal("The function is not working properly.")
				}
			},
		},
	}

	for _, v := range tests {
		t.Run(v.name, v.runner)
	}
}

func Test_FetchNowTime(t *testing.T) {
	tz := "Asia/Tokyo"

	tests := []struct {
		name   string
		runner func(t *testing.T)
	}{
		{
			name: "Check that the function returns a value without error.",
			runner: func(t *testing.T) {
				if err := os.Setenv("TIME_ZONE", tz); err != nil {
					t.Fatalf("Failed to set file TIME_ZONE environment variables.")
				}

				nt, err := FetchNowTime()
				et := time.Time{}
				if nt == et || err != nil {
					t.Fatalf("Failed to fetch now time.")
				}
			},
		},
		{
			name: "Failed to fetch time zone",
			runner: func(t *testing.T) {
				// Unset environment variables to reproduce the condition.
				if err := os.Unsetenv("TIME_ZONE"); err != nil {
					t.Fatalf("Failed to unset file TIME_ZONE environment variables.")
				}

				if _, err := FetchNowTime(); err == nil {
					t.Fatalf("Not behaving as intended.")
				}
			},
		},
		{
			name: "Check that it is unable to fetch now time if setting time zone with mistakes.",
			runner: func(t *testing.T) {
				if err := os.Setenv("TIME_ZONE", "xxx"); err != nil {
					t.Fatalf("Failed to set file TIME_ZONE environment variables.")
				}

				if _, err := FetchNowTime(); err == nil {
					t.Fatalf("Not behaving as intended.")
				}
			},
		},
	}

	for _, v := range tests {
		t.Run(v.name, v.runner)
	}
}
