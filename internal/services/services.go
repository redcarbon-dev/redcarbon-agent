package services

import (
	"context"

	agents_publicv1 "pkg.redcarbon.ai/proto/redcarbon/agents_public/v1"
	"pkg.redcarbon.ai/proto/redcarbon/agents_public/v1/agents_publicv1connect"
)

type Service interface {
	RunService(ctx context.Context)
}

func NewServicesFromConfig(agentsCli agents_publicv1connect.AgentsPublicAPIsV1SrvClient, config *agents_publicv1.AgentConfiguration) []Service {
	services := []Service{}

	if config.GetQradarJobConfiguration() != nil {
		services = append(services, newQRadarService(config.GetQradarJobConfiguration(), agentsCli))
	}

	if config.GetSentineloneJobConfiguration() != nil {
		services = append(services, newSentinelOneService(config.GetSentineloneJobConfiguration(), agentsCli))
	}

	return services
}
