package grayLogImpossibleTravel

import (
	"context"
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

func NewGrayLogImpossibleTravelService(conf *agentsExternalApiV1.AgentConfiguration, cli agentsExternalApiV1.AgentsExternalV1SrvClient) *ServiceGrayLogImpossibleTravel {
	itConf := conf.Data.GetGraylogImpossibleTravel()

	return &ServiceGrayLogImpossibleTravel{
		ac:     conf,
		aCli:   cli,
		itConf: itConf,
		glCli:  graylog.NewGrayLogClient(itConf.Token, itConf.Url, itConf.SkipSsl),
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

	impossibleTravels := s.findImpossibleTravel(logs, s.ac.AgentConfigurationId)

	for _, it := range impossibleTravels {
		_, err := s.aCli.SendGrayLogImpossibleTravelData(ctx, it)
		if err != nil {
			l.Errorf("Error while sending impossible travel for error %v - data %v", err, it)
			continue
		}
	}
}

func (s ServiceGrayLogImpossibleTravel) findImpossibleTravel(logs []map[string]string, acID string) []*agentsExternalApiV1.SendGrayLogImpossibleTravelDataReq {
	var finds []*agentsExternalApiV1.SendGrayLogImpossibleTravelDataReq

	byUser := utils.GroupMapByColumn(logs, "UserId")

	for user, userLogs := range byUser {
		c := utils.GetUniqueDataForColumnInMap(userLogs, "ClientIP_country_code")
		if len(c) <= 1 {
			continue
		}

		var its []*agentsExternalApiV1.GrayLogImpossibleTravelLog

		for _, userLog := range userLogs {
			its = append(its, &agentsExternalApiV1.GrayLogImpossibleTravelLog{
				Logs: userLog,
			})
		}

		ips := utils.GetUniqueDataForColumnInMap(userLogs, "ClientIP")

		finds = append(finds, &agentsExternalApiV1.SendGrayLogImpossibleTravelDataReq{
			Ips:                  ips,
			User:                 user,
			Countries:            c,
			AgentConfigurationId: acID,
			ImpossibleTravelLogs: its,
		})
	}

	return finds
}
