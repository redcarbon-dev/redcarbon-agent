package graylog_datamine

import (
	"context"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/timestamppb"

	"pkg.redcarbon.ai/internal/graylog"
	agentsExternalApiV1 "pkg.redcarbon.ai/proto/redcarbon/external_api/agents/api/v1"
)

type ServiceGrayLogDataMine struct {
	ac      *agentsExternalApiV1.AgentConfiguration
	gdmConf *agentsExternalApiV1.GrayLogDataMineData
	aCli    agentsExternalApiV1.AgentsExternalV1SrvClient
	glCli   graylog.Client
}

func NewGrayLogDataMineService(conf *agentsExternalApiV1.AgentConfiguration, cli agentsExternalApiV1.AgentsExternalV1SrvClient) *ServiceGrayLogDataMine {
	gdmConf := conf.Data.GetGraylogDatamine()

	return &ServiceGrayLogDataMine{
		ac:      conf,
		aCli:    cli,
		gdmConf: gdmConf,
		glCli:   graylog.NewGrayLogClient(gdmConf.Token, gdmConf.Url, gdmConf.SkipSsl),
	}
}

func (s ServiceGrayLogDataMine) RunService(ctx context.Context) {
	l := logrus.WithField("configurationId", s.ac.AgentConfigurationId)

	l.Infof("Retrieving pending queries...")

	qs, err := s.aCli.GetGrayLogDataMinePendingQueries(ctx, &agentsExternalApiV1.GetGrayLogDataMinePendingQueriesReq{
		AgentConfigurationId: s.ac.AgentConfigurationId,
	})
	if err != nil {
		l.Errorf("Error while retrieving the pending queries for error %v", err)
		return
	}

	var wg sync.WaitGroup

	for _, q := range qs.GraylogDatamineQueries {
		wg.Add(1)

		go func(query *agentsExternalApiV1.GrayLogDataMineQuery) {
			defer wg.Done()

			err := s.runQuery(ctx, query, l)
			if err == nil {
				return
			}

			l.Errorf("Error while executing the query %s for error %v", query.Id, err)

			_, err = s.aCli.SendGrayLogDatamineQueryErrorData(ctx, &agentsExternalApiV1.SendGrayLogDatamineQueryErrorDataReq{
				AgentConfigurationId: s.ac.AgentConfigurationId,
				QueryId:              query.Id,
				Error:                err.Error(),
			})
			if err != nil {
				l.Errorf("Error while sending the error due to %v for query %s", err, query.Id)
			}
		}(q)
	}

	wg.Wait()

	l.Infof("Done")
}

func (s ServiceGrayLogDataMine) runQuery(ctx context.Context, q *agentsExternalApiV1.GrayLogDataMineQuery, l *logrus.Entry) error {
	l.Infof("Starting query %s...", q.Id)

	res, err := s.glCli.QueryData(ctx, q.Query, q.SearchStartTime.AsTime(), q.SearchStopTime.AsTime(), []string{"timestamp", "gl2_message_id", "source", "message"})
	if err != nil {
		return err
	}

	var results []*agentsExternalApiV1.GrayLogDataMineResult

	for _, v := range res {
		t, err := time.Parse(time.RFC3339Nano, v["timestamp"])
		if err != nil {
			return err
		}

		results = append(results, &agentsExternalApiV1.GrayLogDataMineResult{
			Uuid:      v["gl2_message_id"],
			Source:    v["source"],
			Message:   v["message"],
			Timestamp: timestamppb.New(t),
		})
	}

	_, err = s.aCli.SendGrayLogDatamineQueryResultsData(ctx, &agentsExternalApiV1.SendGrayLogDatamineQueryResultsDataReq{
		QueryId:              q.Id,
		AgentConfigurationId: s.ac.AgentConfigurationId,
		Results:              results,
	})
	if err != nil {
		return err
	}

	l.Infof("Query %s successfully executed", q.Id)

	return nil
}