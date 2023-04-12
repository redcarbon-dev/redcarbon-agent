package routines

import (
	"github.com/google/go-github/v50/github"
	"pkg.redcarbon.ai/internal/auth"
	agentsExternalApiV1 "pkg.redcarbon.ai/proto/redcarbon/external_api/agents/api/v1"
)

type routineConfig struct {
	agentsCli agentsExternalApiV1.AgentsExternalV1SrvClient
	authSrv   auth.AuthenticationService
	gh        *github.Client
	done      chan bool
}

func NewRoutineJobs(a agentsExternalApiV1.AgentsExternalV1SrvClient, authSrv auth.AuthenticationService, gh *github.Client, done chan bool) routineConfig {
	return routineConfig{
		agentsCli: a,
		authSrv:   authSrv,
		gh:        gh,
		done:      done,
	}
}
