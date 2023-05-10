package services

import (
	"context"

	graylog_datamine "pkg.redcarbon.ai/internal/graylog-datamine"
	"pkg.redcarbon.ai/internal/graylog-impossible-travel"
	"pkg.redcarbon.ai/internal/sentinelone"
	agentsPublicApiV1 "pkg.redcarbon.ai/proto/redcarbon/public_apis/agents/api/v1"
)

type Service interface {
	RunService(ctx context.Context)
}

func NewServiceFromConfiguration(conf *agentsPublicApiV1.AgentConfiguration, cli agentsPublicApiV1.AgentsPublicApiV1SrvClient) Service {
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
