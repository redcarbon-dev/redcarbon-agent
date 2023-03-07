package services

import (
	"context"

	graylog_datamine "pkg.redcarbon.ai/internal/graylog-datamine"
	"pkg.redcarbon.ai/internal/graylog-impossible-travel"
	"pkg.redcarbon.ai/internal/sentinelone"
	agentsExternalApiV1 "pkg.redcarbon.ai/proto/redcarbon/external_api/agents/api/v1"
)

type Service interface {
	RunService(ctx context.Context)
}

func NewServiceFromConfiguration(conf *agentsExternalApiV1.AgentConfiguration, cli agentsExternalApiV1.AgentsExternalV1SrvClient) Service {
	if conf.Data.GetSentinelOne() != nil {
		return sentinelone.NewSentinelOneService(conf, cli)
	}

	if conf.Data.GetGraylogImpossibleTravel() != nil {
		return grayLogImpossibleTravel.NewGrayLogImpossibleTravelService(conf, cli)
	}

	if conf.Data.GetGraylogDatamine() != nil {
		return graylog_datamine.NewGrayLogDataMineService(conf, cli)
	}

	return nil
}
