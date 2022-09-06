package util

import (
	"os"
	"time"
)

func GetLocalTime(t time.Time) time.Time {
	loc, err := time.LoadLocation(os.Getenv("TIME_ZONE"))
	if err == nil {
		t = t.In(loc)
	}

	return t
}
