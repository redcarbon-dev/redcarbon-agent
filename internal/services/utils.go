package services

import (
	"time"
)

func retrieveStartAndEndTime(timeToStart time.Time) (start, end time.Time) {
	end = time.Now()

	start = timeToStart
	if start.IsZero() {
		start = end.Add(-hoursToFetch * time.Hour)
	}

	return
}
