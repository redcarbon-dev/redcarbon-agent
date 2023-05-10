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
	"pkg.redcarbon.ai/internal/services"
	"pkg.redcarbon.ai/mocks"
	agentsPublicApiV1 "pkg.redcarbon.ai/proto/redcarbon/public_apis/agents/api/v1"
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
	}))

	defer ts.Close()

	f, err := os.CreateTemp("", "rcagent-sentinelone-test")
	assert.Nil(t, err)

	defer f.Close()

	viper.SetConfigFile(f.Name())
	viper.SetConfigType("yaml")

	cli := mocks.AgentsPublicApiV1SrvClient{}

	cli.On("SendSentinelOneData", mock.Anything, mock.Anything).Return(&agentsPublicApiV1.SendSentinelOneDataRes{ReceivedAt: timestamppb.Now()}, nil)

	s := services.NewServiceFromConfiguration(&agentsPublicApiV1.AgentConfiguration{
		AgentConfigurationId: "cf:1234567890",
		Name:                 "test",
		Type:                 "sentinel_one",
		CreatedAt:            timestamppb.Now(),
		UpdatedAt:            timestamppb.Now(),
		Data: &agentsPublicApiV1.AgentConfigurationData{
			Data: &agentsPublicApiV1.AgentConfigurationData_SentinelOne{
				SentinelOne: &agentsPublicApiV1.SentinelOneData{
					Url:      ts.URL,
					ApiToken: "xxx",
				},
			},
		},
	}, &cli)

	s.RunService(context.Background())

	cli.AssertNumberOfCalls(t, "SendSentinelOneData", 2)
	assert.NotEqual(t, time.Time{}, viper.Get("configurations.cf:1234567890.from"))
}
