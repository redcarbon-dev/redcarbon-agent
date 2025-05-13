package services

import (
	"context"
	"encoding/json"
	"fmt"

	"connectrpc.com/connect"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"pkg.redcarbon.ai/internal/config"
	"pkg.redcarbon.ai/internal/services/fortisiem"
	agents_publicv1 "pkg.redcarbon.ai/proto/redcarbon/agents_public/v1"
	"pkg.redcarbon.ai/proto/redcarbon/agents_public/v1/agents_publicv1connect"
)

type srvFortiSIEM struct {
	cli       fortisiem.FortiSIEMClient
	agentsCli agents_publicv1connect.AgentsPublicAPIsV1SrvClient
	profile   config.Profile
}

func newFortiSIEMService(conf *agents_publicv1.FortiSIEMJobConfiguration, agentsCli agents_publicv1connect.AgentsPublicAPIsV1SrvClient, p config.Profile) Service {
	return &srvFortiSIEM{
		cli:       fortisiem.NewFortiSIEMClient(conf.Host, conf.Username, conf.Password, conf.VerifySsl),
		agentsCli: agentsCli,
		profile:   p,
	}
}

func (s srvFortiSIEM) RunService(ctx context.Context) {
	l := logrus.WithFields(logrus.Fields{
		"service": "fortisiem",
		"trace":   uuid.NewString(),
		"profile": s.profile.Name,
	})

	l.Info("Starting FortiSIEM service")
	s.profile = config.LoadProfile(s.profile.Name)

	start, end := retrieveStartAndEndTime(s.profile.FortiSIEM.LastExecution)

	alerts, err := s.cli.FetchAlerts(ctx, start, end)
	if err != nil {
		l.WithError(err).Error("Error while fetching the alerts")
		return
	}

	l.Infof("Found %d alerts", len(alerts))

	for _, alert := range alerts {
		incident, err := s.buildIncidentToIngest(alert)
		if err != nil {
			l.WithError(err).Warn("Error while building the incident to ingest for alert")
		}

		req := connect.NewRequest(incident)
		req.Header().Set("authorization", fmt.Sprintf("ApiToken %s", s.profile.Profile.Token))

		_, err = s.agentsCli.IngestIncident(ctx, req)
		if err != nil {
			l.WithError(err).Error("failed to ingest incident")
			continue
		}
	}

	s.profile.FortiSIEM.LastExecution = end

	config.OverwriteProfileInConfig(l, s.profile)

	l.Info("FortiSIEM service completed")
}

func (s srvFortiSIEM) buildIncidentToIngest(incident map[string]interface{}) (*agents_publicv1.IngestIncidentRequest, error) {
	iStr, err := json.Marshal(incident)
	if err != nil {
		return nil, err
	}

	var inc fortisiem.Incident

	err = json.Unmarshal(iStr, &inc)
	if err != nil {
		return nil, err
	}

	idStr := fmt.Sprintf("%d", inc.IncidentID)

	return &agents_publicv1.IngestIncidentRequest{
		Title:       inc.IncidentTitle,
		Description: inc.IncidentTitle,
		RawData:     string(iStr),
		Severity:    mapFortiSIEMSeverity(inc.EventSeverity),
		Origin:      "fortisiem",
		OriginalId:  &idStr,
		OriginalUrl: nil,
	}, nil
}

func mapFortiSIEMSeverity(sev int) uint32 {
	if sev >= 1 && sev <= 4 {
		return uint32(10)
	}

	if sev >= 5 && sev <= 8 {
		return uint32(40)
	}

	return uint32(70)
}
