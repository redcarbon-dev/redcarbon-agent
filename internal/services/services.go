package services

import (
	"context"
	"pkg.redcarbon.ai/internal/config"

	agents_publicv1 "pkg.redcarbon.ai/proto/redcarbon/agents_public/v1"
	"pkg.redcarbon.ai/proto/redcarbon/agents_public/v1/agents_publicv1connect"
)

type Service interface {
	RunService(ctx context.Context)
}

func NewServicesFromConfig(agentsCli agents_publicv1connect.AgentsPublicAPIsV1SrvClient, conf *agents_publicv1.AgentConfiguration, p config.Profile) []Service {
	services := []Service{}

	if p.Debug.Active {
		return []Service{newDebugService(p)}
	}

	if conf.GetQradarJobConfiguration() != nil {
		services = append(services, newQRadarService(conf.GetQradarJobConfiguration(), agentsCli, p))
	}

	if conf.GetSentineloneJobConfiguration() != nil {
		services = append(services, newSentinelOneService(conf.GetSentineloneJobConfiguration(), agentsCli, p))
	}

	if conf.GetFortisiemJobConfiguration() != nil {
		services = append(services, newFortiSIEMService(conf.GetFortisiemJobConfiguration(), agentsCli, p))
	}

	return services
}
