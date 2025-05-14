package services

import (
	"context"
	"encoding/json"
	"fmt"

	"connectrpc.com/connect"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"pkg.redcarbon.ai/internal/config"
	"pkg.redcarbon.ai/internal/services/sentinel_one"
	agents_publicv1 "pkg.redcarbon.ai/proto/redcarbon/agents_public/v1"
	"pkg.redcarbon.ai/proto/redcarbon/agents_public/v1/agents_publicv1connect"
)

type srvSentinel struct {
	cli       sentinel_one.SentinelOneClient
	agentsCli agents_publicv1connect.AgentsPublicAPIsV1SrvClient
	profile   config.Profile
}

func newSentinelOneService(conf *agents_publicv1.SentinelOneJobConfiguration, agentsCli agents_publicv1connect.AgentsPublicAPIsV1SrvClient, p config.Profile) Service {
	return &srvSentinel{
		cli:       sentinel_one.NewSentinelOneClient(conf.Host, conf.Token, conf.VerifySsl),
		agentsCli: agentsCli,
		profile:   p,
	}
}

func (s srvSentinel) RunService(ctx context.Context) {
	l := logrus.WithFields(logrus.Fields{
		"service": "sentinel-one",
		"trace":   uuid.NewString(),
		"profile": s.profile.Name,
	})

	l.Info("Starting SentinelOne service")
	s.profile = config.LoadProfile(s.profile.Name)

	start, end := retrieveStartAndEndTime(s.profile.SentinelONE.LastExecution)

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
		req.Header().Set("authorization", fmt.Sprintf("ApiToken %s", s.profile.Profile.Token))

		_, err = s.agentsCli.IngestIncident(ctx, req)
		if err != nil {
			l.WithError(err).Error("failed to ingest incident")
			continue
		}
	}

	s.profile.SentinelONE.LastExecution = end

	config.OverwriteProfileInConfig(l, s.profile)

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
