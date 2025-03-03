package services

import (
	"context"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
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
	l := logrus.WithFields(logrus.Fields{
		"service": "debug",
		"trace":   uuid.NewString(),
	})

	s.profile = config.LoadProfile(s.profile.Name)

	l.Info("Starting Debug service")
	start, end := retrieveStartAndEndTime(s.profile.Debug.LastExecution)

	l.Infof("Debug from %s to %s", start, end)

	s.profile.Debug.LastExecution = end

	config.OverwriteProfileInConfig(l, s.profile)

	l.Info("Debug service completed")
}
