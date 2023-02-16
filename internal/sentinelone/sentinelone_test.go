package sentinelone_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"

	"pkg.redcarbon.ai/internal/sentinelone"
	"pkg.redcarbon.ai/mocks"
	agentsExternalApiV1 "pkg.redcarbon.ai/proto/redcarbon/external_api/agents/api/v1"
)

var s *grpc.Server

func TestShouldSendAllTheSentinelOneData(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		assert.Equal(t, "ApiToken xxx", r.Header.Get("Authorization"))

		query := r.URL.Query()
		limit := query.Get("limit")
		skip := query.Get("skip")

		assert.Equal(t, limit, "100")

		if skip == "0" {
			res, err := json.Marshal(sentinelone.FetchResponse{Data: []map[string]interface{}{
				{"Hello": "World"},
				{"bar": []string{"foo"}},
			}})

			assert.Nil(t, err)

			_, err = w.Write(res)
			assert.Nil(t, err)
		} else {
			res, err := json.Marshal([]string{})

			assert.Nil(t, err)

			_, err = w.Write(res)
			assert.Nil(t, err)
		}

		w.WriteHeader(http.StatusOK)
	}))

	defer ts.Close()

	f, err := os.CreateTemp("", "rcagent-sentinelone-test")
	assert.Nil(t, err)

	defer f.Close()

	viper.SetConfigFile(f.Name())
	viper.SetConfigType("yaml")

	cli := mocks.AgentsExternalV1SrvClient{}

	cli.On("SendData", mock.Anything, mock.Anything).Return(&agentsExternalApiV1.SendDataRes{ReceivedAt: timestamppb.Now()}, nil)

	sentinelone.RunSentinelOneService(context.Background(), &agentsExternalApiV1.AgentConfiguration{
		AgentConfigurationId: "cf:1234567890",
		Name:                 "test",
		Type:                 "sentinel_one",
		CreatedAt:            timestamppb.Now(),
		UpdatedAt:            timestamppb.Now(),
		Data: &agentsExternalApiV1.AgentConfigurationData{
			Data: &agentsExternalApiV1.AgentConfigurationData_SentinelOne{
				SentinelOne: &agentsExternalApiV1.SentinelOneData{
					Url:      ts.URL,
					ApiToken: "xxx",
				},
			},
		},
	}, &cli)

	cli.AssertNumberOfCalls(t, "SendData", 2)
	assert.NotEqual(t, time.Time{}, viper.Get("configurations.cf:1234567890.from"))
}
