package services

import (
	"connectrpc.com/connect"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"pkg.redcarbon.ai/internal/services/qradar"
	agents_publicv1 "pkg.redcarbon.ai/proto/redcarbon/agents_public/v1"
	"pkg.redcarbon.ai/proto/redcarbon/agents_public/v1/agents_publicv1connect"
	"strings"
	"time"
)

type srv struct {
	cli       qradar.QRadarClient
	agentsCli agents_publicv1connect.AgentsPublicAPIsV1SrvClient
}

type incidentToIngest struct {
	Incident                  map[string]interface{}   `json:"incident"`
	OffenseType               map[string]interface{}   `json:"offense_type"`
	LocalDestinationAddresses []map[string]interface{} `json:"local_destination_addresses"`
	SourceAddresses           []map[string]interface{} `json:"source_addresses"`
}

func NewService(conf *agents_publicv1.QRadarJobConfiguration, agentsCli agents_publicv1connect.AgentsPublicAPIsV1SrvClient) Service {
	return &srv{
		cli:       qradar.NewQRadarClient(conf.Host, conf.Token, conf.VerifySsl),
		agentsCli: agentsCli,
	}
}

func (s srv) RunService(ctx context.Context) {
	l := logrus.WithFields(logrus.Fields{
		"service": "qradar",
		"trace":   uuid.NewString(),
	})

	l.Info("Starting QRadar service")
	start, end := retrieveTimeRangeForSearch()
	offenses, err := s.cli.FetchOffenses(ctx, start, end)
	if err != nil {
		return
	}

	l.Infof("Found %d offenses", len(offenses))

	for _, offense := range offenses {
		i, err := s.retrieveIncidentDataForOffense(ctx, offense)
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

	viper.Set("qradar.last_execution", end)
	if err := viper.WriteConfig(); err != nil {
		l.WithError(err).Error("failed to write config")
	}

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

	i := incidentToIngest{
		Incident:                  offense,
		OffenseType:               ot,
		LocalDestinationAddresses: la,
		SourceAddresses:           sa,
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
	if magnitude <= 6 {
		return 10
	}

	if magnitude >= 7 && magnitude <= 8 {
		return 40
	}

	if magnitude == 9 {
		return 70
	}

	return 90
}

func retrieveTimeRangeForSearch() (start, end time.Time) {
	end = time.Now()
	start = viper.GetTime("qradar.last_execution")

	if start.IsZero() {
		start = end.Add(-4 * time.Hour)
	}

	return
}
