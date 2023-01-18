package routines

import agentsExternalApiV1 "pkg.redcarbon.ai/proto/redcarbon/external_api/agents/api/v1"

type routineConfig struct {
	agentsCli agentsExternalApiV1.AgentsExternalV1SrvClient
}

func NewRoutineJobs(a agentsExternalApiV1.AgentsExternalV1SrvClient) routineConfig {
	return routineConfig{
		agentsCli: a,
	}
}
