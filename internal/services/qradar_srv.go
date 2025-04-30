package services

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"connectrpc.com/connect"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"pkg.redcarbon.ai/internal/config"
	"pkg.redcarbon.ai/internal/services/qradar"
	agents_publicv1 "pkg.redcarbon.ai/proto/redcarbon/agents_public/v1"
	"pkg.redcarbon.ai/proto/redcarbon/agents_public/v1/agents_publicv1connect"
)

const (
	maximumMagnitudeLow    = 6
	minimumMagnitudeMedium = 7
	maximumMagnitudeMedium = 8
	magnitudeHigh          = 9

	rcSeverityLow      = 10
	rcSeverityMedium   = 40
	rcSeverityHigh     = 70
	rcSeverityCritical = 90

	hoursToFetch = 2
)

type srv struct {
	cli       qradar.QRadarClient
	agentsCli agents_publicv1connect.AgentsPublicAPIsV1SrvClient
	profile   config.Profile
}

type incidentToIngest struct {
	Incident                  map[string]any   `json:"incident"`
	OffenseType               map[string]any   `json:"offense_type"`
	LocalDestinationAddresses []map[string]any `json:"local_destination_addresses"`
	SourceAddresses           []map[string]any `json:"source_addresses"`
	Events                    []map[string]any `json:"events"`
}

func newQRadarService(conf *agents_publicv1.QRadarJobConfiguration, agentsCli agents_publicv1connect.AgentsPublicAPIsV1SrvClient, p config.Profile) Service {
	return &srv{
		cli:       qradar.NewQRadarClient(conf.Host, conf.Token, conf.VerifySsl),
		agentsCli: agentsCli,
		profile:   p,
	}
}

func (s srv) RunService(ctx context.Context) {
	l := logrus.WithFields(logrus.Fields{
		"service": "qradar",
		"trace":   uuid.NewString(),
	})

	l.Info("Starting QRadar service")
	s.profile = config.LoadProfile(s.profile.Name)

	start, end := retrieveStartAndEndTime(s.profile.QRadar.LastExecution)

	offenses, err := s.cli.FetchOffenses(ctx, start, end)
	if err != nil {
		l.WithError(err).Error("failed to fetch offenses")
		return
	}

	l.Infof("Found %d offenses", len(offenses))

	for _, offense := range offenses {
		i, err := s.retrieveIncidentDataForOffense(ctx, offense)
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

	s.profile.QRadar.LastExecution = end

	config.OverwriteProfileInConfig(l, s.profile)

	l.Info("QRadar service completed")
}

func (s srv) retrieveIncidentDataForOffense(ctx context.Context, offense map[string]interface{}) (*agents_publicv1.IngestIncidentRequest, error) {
	of, err := parseOffense(offense)
	if err != nil {
		return nil, err
	}

	ot, err := s.cli.FetchOffenseType(ctx, of.OffenseType)
	if err != nil {
		return nil, err
	}

	la, err := s.SafeFetchIds(ctx, of.LocalDestinationAddressIds, s.cli.FetchLocalDestinationAddresses)
	if err != nil {
		return nil, err
	}

	sa, err := s.SafeFetchIds(ctx, of.SourceAddressIds, s.cli.FetchSourceAddresses)
	if err != nil {
		return nil, err
	}

	events := s.cli.SafeSearchOffenseEvents(ctx, of.ID, int64(of.StartTime))

	i := incidentToIngest{
		Incident:                  offense,
		OffenseType:               ot,
		LocalDestinationAddresses: la,
		SourceAddresses:           sa,
		Events:                    events,
	}

	iStr, err := json.Marshal(i)
	if err != nil {
		return nil, err
	}

	idStr := fmt.Sprintf("%d", of.ID)

	return &agents_publicv1.IngestIncidentRequest{
		Title:       strings.TrimSpace(of.Description),
		Description: of.Description,
		RawData:     string(iStr),
		Severity:    mapSeverity(of.Magnitude),
		Origin:      "qradar",
		OriginalId:  &idStr,
		OriginalUrl: s.cli.RetrieveOffenseUrl(of.ID),
	}, nil
}

func parseOffense(offense map[string]interface{}) (qradar.Offense, error) {
	ofStr, err := json.Marshal(offense)
	if err != nil {
		return qradar.Offense{}, err
	}

	var of qradar.Offense

	err = json.Unmarshal(ofStr, &of)
	if err != nil {
		return qradar.Offense{}, err
	}

	return of, nil
}

func (s srv) SafeFetchIds(ctx context.Context, ids []int, f func(context.Context, []int) ([]map[string]interface{}, error)) ([]map[string]interface{}, error) {
	if len(ids) == 0 {
		return []map[string]interface{}{}, nil
	}

	return f(ctx, ids)
}

func mapSeverity(magnitude int) uint32 {
	if magnitude <= maximumMagnitudeLow {
		return rcSeverityLow
	}

	if magnitude >= minimumMagnitudeMedium && magnitude <= maximumMagnitudeMedium {
		return rcSeverityMedium
	}

	if magnitude == magnitudeHigh {
		return rcSeverityHigh
	}

	return rcSeverityCritical
}
