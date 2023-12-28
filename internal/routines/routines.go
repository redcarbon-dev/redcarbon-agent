package routines

import (
	"github.com/google/go-github/v50/github"
	"pkg.redcarbon.ai/proto/redcarbon/agents_public/v1/agents_publicv1connect"
)

type RoutineConfig struct {
	agentsCli agents_publicv1connect.AgentsPublicAPIsV1SrvClient
	gh        *github.Client
	done      chan bool
}

func NewRoutineJobs(a agents_publicv1connect.AgentsPublicAPIsV1SrvClient, gh *github.Client, done chan bool) RoutineConfig {
	return RoutineConfig{
		agentsCli: a,
		gh:        gh,
		done:      done,
	}
}
