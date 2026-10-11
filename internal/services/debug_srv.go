package services

import (
	"context"

	"pkg.redcarbon.ai/internal/config"
)

type srvDebug struct {
	profile config.Profile
}

func newDebugService(p config.Profile) Service {
	return &srvDebug{
		profile: p,
	}
}

func (s srvDebug) RunService(ctx context.Context) {
	l := serviceLogger("debug", s.profile)

	s.profile = config.LoadProfile(s.profile.Name)

	l.Info("Starting Debug service")
	start, end := retrieveStartAndEndTime(s.profile.Debug.LastExecution)

	l.Info("Debug window", "start", start, "end", end)

	s.profile.Debug.LastExecution = end

	config.OverwriteProfileInConfig(l, s.profile)

	l.Info("Debug service completed")
}
