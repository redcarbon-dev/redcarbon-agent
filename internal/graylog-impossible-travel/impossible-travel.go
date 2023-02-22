package grayLogImpossibleTravel

import (
	"context"
	"encoding/json"
	"time"

	"github.com/sirupsen/logrus"

	"pkg.redcarbon.ai/internal/graylog"
	"pkg.redcarbon.ai/internal/utils"
	agentsExternalApiV1 "pkg.redcarbon.ai/proto/redcarbon/external_api/agents/api/v1"
)

type ServiceGrayLogImpossibleTravel struct {
	ac     *agentsExternalApiV1.AgentConfiguration
	itConf *agentsExternalApiV1.GrayLogImpossibleTravelData
	aCli   agentsExternalApiV1.AgentsExternalV1SrvClient
	glCli  graylog.Client
}

type findings struct {
	User      string              `json:"user"`
	Ips       []string            `json:"ips"`
	Countries []string            `json:"countries"`
	Logs      []map[string]string `json:"logs"`
}

func NewGrayLogImpossibleTravelService(conf *agentsExternalApiV1.AgentConfiguration, cli agentsExternalApiV1.AgentsExternalV1SrvClient) *ServiceGrayLogImpossibleTravel {
	itConf := conf.Data.GetGraylogImpossibleTravel()

	return &ServiceGrayLogImpossibleTravel{
		ac:     conf,
		aCli:   cli,
		itConf: itConf,
		glCli:  graylog.NewGrayLogClient(itConf.Token, itConf.Url, false),
	}
}

func (s ServiceGrayLogImpossibleTravel) RunService(ctx context.Context) {
	l := logrus.WithField("configurationId", s.ac.AgentConfigurationId)

	to := time.Now()
	from := to.Add(-s.itConf.TimeWindow.AsDuration())

	logs, _ := s.glCli.QueryData(
		ctx,
		"Workload:AzureActiveDirectory AND Operation:UserLoggedIn AND NOT _exists_:customer_subnet_lookup_match",
		from,
		to,
		[]string{
			"timestamp",
			"source",
			"UserId",
			"ClientIP",
			"ClientIP_country_code",
			"UserAgent",
		},
	)

	finds := s.findImpossibleTravel(logs)

	for _, find := range finds {
		data, err := json.Marshal(find)
		if err != nil {
			l.Errorf("Error while converting the data in json for error %v - data %v", err, find)
			continue
		}

		_, err = s.aCli.SendData(ctx, &agentsExternalApiV1.SendDataReq{
			Data:                 string(data),
			DataType:             agentsExternalApiV1.DataType_GRAYLOG_IMPOSSIBLE_TRAVEL,
			AgentConfigurationId: s.ac.AgentConfigurationId,
		})
		if err != nil {
			l.Errorf("Error while sending data for error %v - data %s", err, string(data))
			continue
		}
	}
}

func (s ServiceGrayLogImpossibleTravel) findImpossibleTravel(logs []map[string]string) []findings {
	var finds []findings

	byUser := utils.GroupMapByColumn(logs, "UserId")

	for user, userLogs := range byUser {
		c := utils.GetUniqueDataForColumnInMap(userLogs, "ClientIP_country_code")
		if len(c) <= 1 {
			continue
		}

		ips := utils.GetUniqueDataForColumnInMap(userLogs, "ClientIP")

		finds = append(finds, findings{
			Ips:       ips,
			Logs:      userLogs,
			User:      user,
			Countries: c,
		})
	}

	return finds
}
