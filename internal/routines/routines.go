package routines

import (
	"log/slog"

	"pkg.redcarbon.ai/internal/cli"
	"pkg.redcarbon.ai/internal/config"
	"pkg.redcarbon.ai/internal/github"
	"pkg.redcarbon.ai/proto/redcarbon/agents_public/v1/agents_publicv1connect"
)

type RoutineConfig struct {
	agentsCli agents_publicv1connect.AgentsPublicAPIsV1SrvClient
	gh        *github.Client
	done      chan bool
	profile   config.Profile
}

func NewRoutineJobs(profile config.Profile, clientFactory cli.ClientFactory, gh *github.Client, done chan bool) RoutineConfig {
	return RoutineConfig{
		profile:   profile,
		agentsCli: clientFactory.GetAgentClientForProfile(profile),
		gh:        gh,
		done:      done,
	}
}

func (r RoutineConfig) Logger() *slog.Logger {
	return slog.With("profile", r.profile.Name)
}
