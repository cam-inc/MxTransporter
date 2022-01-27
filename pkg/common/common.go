package common

import (
	"mxtransporter/config"
	"mxtransporter/pkg/errors"
	"time"
)

// Contains check if the array has the specified string
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
		return time.Time{}, errors.InternalServerError.Wrap("Failed to fetch time zone.", err)
	}

	tl, err := time.LoadLocation(tz)
	if err != nil {
		return time.Time{}, errors.InternalServerError.Wrap("The timezone could not be converted to location type.", err)
	}

	return time.Now().In(tl), nil
}
