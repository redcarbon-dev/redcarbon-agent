package services

import (
	"connectrpc.com/connect"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"pkg.redcarbon.ai/internal/services/sentinel_one"
	agents_publicv1 "pkg.redcarbon.ai/proto/redcarbon/agents_public/v1"
	"pkg.redcarbon.ai/proto/redcarbon/agents_public/v1/agents_publicv1connect"
)

type srvSentinel struct {
	cli       sentinel_one.SentinelOneClient
	agentsCli agents_publicv1connect.AgentsPublicAPIsV1SrvClient
}

func newSentinelOneService(conf *agents_publicv1.SentinelOneJobConfiguration, agentsCli agents_publicv1connect.AgentsPublicAPIsV1SrvClient) Service {
	return &srvSentinel{
		cli:       sentinel_one.NewSentinelOneClient(conf.Host, conf.Token, conf.VerifySsl),
		agentsCli: agentsCli,
	}
}

func (s srvSentinel) RunService(ctx context.Context) {
	l := logrus.WithFields(logrus.Fields{
		"service": "sentinel-one",
		"trace":   uuid.NewString(),
	})

	l.Info("Starting SentinelOne service")
	start, end := retrieveSearchTimeRangeForKey("sentinel_one")

	threats, err := s.cli.FetchThreats(ctx, start, end)
	if err != nil {
		l.WithError(err).Error("Error while fetching the threats")
		return
	}

	for _, threat := range threats {
		i, err := retrieveIncidentDataForThreat(threat)
		if err != nil {
			continue
		}

		req := connect.NewRequest(i)
		req.Header().Set("authorization", fmt.Sprintf("ApiToken %s", viper.Get("auth.access_token")))

		_, err = s.agentsCli.IngestIncident(ctx, req)
		if err != nil {
			l.WithError(err).Error("failed to ingest incident")
			continue
		}
	}

	viper.Set("sentinel_one.last_execution", end)

	if err := viper.WriteConfig(); err != nil {
		l.WithError(err).Error("failed to write config")
	}

	l.Info("SentinelOne service completed")
}

func retrieveIncidentDataForThreat(threat map[string]interface{}) (*agents_publicv1.IngestIncidentRequest, error) {
	thStr, err := json.Marshal(threat)
	if err != nil {
		return nil, err
	}

	var t sentinel_one.Threat

	err = json.Unmarshal(thStr, &t)
	if err != nil {
		return nil, err
	}

	title := fmt.Sprintf("SentinelOne - %s", t.ThreatInfo.Classification)

	return &agents_publicv1.IngestIncidentRequest{
		Title:       title,
		Description: title,
		RawData:     string(thStr),
		Severity:    0,
		Origin:      "sentinel-one",
		OriginalId:  &t.Id,
	}, nil
}
