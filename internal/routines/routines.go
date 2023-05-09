package routines

import (
	"github.com/google/go-github/v50/github"

	"pkg.redcarbon.ai/internal/auth"
	agentsPublicApiV1 "pkg.redcarbon.ai/proto/redcarbon/public_apis/agents/api/v1"
)

type routineConfig struct {
	agentsCli agentsPublicApiV1.AgentsPublicApiV1SrvClient
	authSrv   auth.AuthenticationService
	gh        *github.Client
	done      chan bool
}

func NewRoutineJobs(a agentsPublicApiV1.AgentsPublicApiV1SrvClient, authSrv auth.AuthenticationService, gh *github.Client, done chan bool) routineConfig {
	return routineConfig{
		agentsCli: a,
		authSrv:   authSrv,
		gh:        gh,
		done:      done,
	}
}
