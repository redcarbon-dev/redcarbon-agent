package routines

import (
	"pkg.redcarbon.ai/internal/auth"
	agentsExternalApiV1 "pkg.redcarbon.ai/proto/redcarbon/external_api/agents/api/v1"
)

type routineConfig struct {
	agentsCli agentsExternalApiV1.AgentsExternalV1SrvClient
	authSrv   auth.AuthenticationService
}

func NewRoutineJobs(a agentsExternalApiV1.AgentsExternalV1SrvClient, authSrv auth.AuthenticationService) routineConfig {
	return routineConfig{
		agentsCli: a,
		authSrv:   authSrv,
	}
}
