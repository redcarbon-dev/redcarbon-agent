package services

import (
	"fmt"
	"github.com/spf13/viper"
	"time"
)

func retrieveSearchTimeRangeForKey(key string) (start, end time.Time) {
	end = time.Now()
	start = viper.GetTime(fmt.Sprintf("%s.last_execution", key))

	if start.IsZero() {
		start = end.Add(-hoursToFetch * time.Hour)
	}

	return
}
