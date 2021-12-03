package common

import (
	"mxtransporter/config"
	"time"
)

// Check if the array has the specified string
func Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func FetchNowTime() (time.Time, error) {
	tz, err := config.FetchTimeZone()
	if err != nil {
		return time.Time{}, err
	}

	tl, err := time.LoadLocation(tz)
	if err != nil {
		return time.Time{}, err
	}

	return time.Now().In(tl), nil
}
