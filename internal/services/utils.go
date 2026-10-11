package services

import (
	"log/slog"
	"time"
	"uuid"

	"pkg.redcarbon.ai/internal/config"
)

func serviceLogger(service string, p config.Profile) *slog.Logger {
	return slog.With(
		"service", service,
		"trace", uuid.New().String(),
		"profile", p.Name,
	)
}

func retrieveStartAndEndTime(timeToStart time.Time) (start, end time.Time) {
	end = time.Now()

	start = timeToStart
	if start.IsZero() {
		start = end.Add(-hoursToFetch * time.Hour)
	}

	return
}
