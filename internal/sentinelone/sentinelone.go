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

func RunSentinelOneService(ctx context.Context, ac *agentsExternalApiV1.AgentConfiguration, aCli agentsExternalApiV1.AgentsExternalV1SrvClient) {
	l := logrus.WithField("configurationId", ac.AgentConfigurationId)

	l.Infof("Starting SentinelOne Configuration...")

	sentinelOneConfig := ac.Data.GetSentinelOne()

	baseUrl, err := url.Parse(sentinelOneConfig.Url)
	if err != nil {
		return
	}

	reqUrl := baseUrl.JoinPath("/web/api/v2.1/threats")

	from, to := retrieveTimeBoundariesForConfiguration(ac)

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

		var parsedBody []any

		err = json.Unmarshal(body, &parsedBody)
		if err != nil {
			return
		}

		if len(parsedBody) == 0 {
			break
		}

		for _, data := range parsedBody {
			l.Infof("Sending %v\n", data)

			dataJ, err := json.Marshal(data)
			if err != nil {
				return
			}

			_, err = aCli.SendData(ctx, &agentsExternalApiV1.SendDataReq{
				Data: string(dataJ),
			})
			if err != nil {
				return
			}
		}

		skip += limit
	}

	viper.Set(fmt.Sprintf("configurations.%s.from", ac.AgentConfigurationId), to)

	err = viper.WriteConfig()
	if err != nil {
		return
	}

	l.Infof("SentinelOne Configuration - Successfully Executed\n")
}

func retrieveTimeBoundariesForConfiguration(ac *agentsExternalApiV1.AgentConfiguration) (time.Time, time.Time) {
	now := time.Now()

	from := viper.GetTime(fmt.Sprintf("configurations.%s.from", ac.AgentConfigurationId))

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
