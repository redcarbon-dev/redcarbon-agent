package sentinelone

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	agentsExternalApiV1 "pkg.redcarbon.ai/proto/redcarbon/external_api/agents/api/v1"
)

const maxDifferenceOfTimes = time.Hour * 24 * 7

type FetchResponse struct {
	Data []map[string]interface{} `json:"data"`
}

type ServiceSentinelOne struct {
	ac   *agentsExternalApiV1.AgentConfiguration
	aCli agentsExternalApiV1.AgentsExternalV1SrvClient
}

func NewSentinelOneService(conf *agentsExternalApiV1.AgentConfiguration, cli agentsExternalApiV1.AgentsExternalV1SrvClient) *ServiceSentinelOne {
	return &ServiceSentinelOne{
		ac:   conf,
		aCli: cli,
	}
}

func (s ServiceSentinelOne) RunService(ctx context.Context) {
	l := logrus.WithField("configurationId", s.ac.AgentConfigurationId)

	l.Infof("Starting SentinelOne Configuration...")

	sentinelOneConfig := s.ac.Data.GetSentinelOne()

	baseUrl, err := url.Parse(sentinelOneConfig.Url)
	if err != nil {
		return
	}

	reqUrl := baseUrl.JoinPath("/web/api/v2.1/threats")

	from, to := s.retrieveTimeBoundariesForConfiguration()

	query := reqUrl.Query()

	query.Set("createdAt__gt", from.Format(time.RFC3339))
	query.Set("createdAt__lt", to.Format(time.RFC3339))

	limit := 100
	skip := 0

	for true {
		query.Set("limit", fmt.Sprintf("%d", limit))
		query.Set("skip", fmt.Sprintf("%d", skip))

		reqUrl.RawQuery = query.Encode()

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqUrl.String(), nil)
		if err != nil {
			return
		}

		req.Header.Add("Authorization", fmt.Sprintf("ApiToken %s", sentinelOneConfig.ApiToken))
		req.Header.Add("User-Agent", "RedCarbon Ingestion Agent")
		req.Header.Add("Accept", "application/json")

		client := http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}

		res, err := client.Do(req)
		if err != nil {
			return
		}

		body, err := io.ReadAll(res.Body)
		if err != nil {
			return
		}

		var parsedBody FetchResponse

		err = json.Unmarshal(body, &parsedBody)
		if err != nil {
			return
		}

		if len(parsedBody.Data) == 0 {
			break
		}

		for _, data := range parsedBody.Data {
			l.Infof("Sending a message...")

			dataJ, err := json.Marshal(data)
			if err != nil {
				logrus.Errorf("Error while marshaling the message %v", err)
				return
			}

			_, err = s.aCli.SendSentinelOneData(ctx, &agentsExternalApiV1.SendSentinelOneDataReq{
				Data:                 string(dataJ),
				AgentConfigurationId: s.ac.AgentConfigurationId,
			})
			if err != nil {
				logrus.Errorf("Error while sending the message %v", err)
				return
			}
		}

		skip += limit
	}

	viper.Set(fmt.Sprintf("configurations.%s.from", s.ac.AgentConfigurationId), to)

	err = viper.WriteConfig()
	if err != nil {
		return
	}

	l.Infof("SentinelOne Configuration - Successfully Executed\n")
}

func (s ServiceSentinelOne) retrieveTimeBoundariesForConfiguration() (time.Time, time.Time) {
	now := time.Now()

	from := viper.GetTime(fmt.Sprintf("configurations.%s.from", s.ac.AgentConfigurationId))

	zero := time.Time{}
	if from == zero {
		return now.Add(-time.Hour), now
	}

	diff := now.Sub(from)
	if diff > maxDifferenceOfTimes {
		return now.Add(-maxDifferenceOfTimes), now
	}

	return from, now
}
