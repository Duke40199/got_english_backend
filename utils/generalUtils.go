package utils

import (
	"time"
)

// GetCurrentTime function is used to get the current time in milliseconds.
func GetCurrentEpochTimeInMiliseconds() int64 {
	var now = time.Now()
	ts := now.UnixNano() / 1000000
	return ts
}
