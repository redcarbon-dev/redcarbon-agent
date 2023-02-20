package impossibleTravel_test

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gocarina/gocsv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"pkg.redcarbon.ai/internal/services"
	"pkg.redcarbon.ai/mocks"
	agentsExternalApiV1 "pkg.redcarbon.ai/proto/redcarbon/external_api/agents/api/v1"
)

type findings struct {
	User      string
	Ips       []string
	Countries []string
	Logs      []map[string]string
}

type graylogData struct {
	ClientIP            string `csv:"ClientIP"`
	ClientIPCountryCode string `csv:"ClientIP_country_code"`
	UserId              string `csv:"UserId"`
}

func TestImpossibleTravel(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		dToken, err := base64.StdEncoding.DecodeString(strings.Replace(r.Header.Get("Authorization"), "Basic ", "", -1))
		assert.Nil(t, err)
		assert.Equal(t, "xxx:token", string(dToken))

		data, err := gocsv.MarshalString([]graylogData{
			{
				ClientIP:            "8.8.8.8",
				ClientIPCountryCode: "US",
				UserId:              "foo",
			},
			{
				ClientIP:            "7.7.7.7",
				ClientIPCountryCode: "US",
				UserId:              "bar",
			},
			{
				ClientIP:            "9.9.9.9",
				ClientIPCountryCode: "DE",
				UserId:              "woo",
			},
			{
				ClientIP:            "9.9.9.9",
				ClientIPCountryCode: "DE",
				UserId:              "woo",
			},
			{
				ClientIP:            "7.7.7.7",
				ClientIPCountryCode: "FR",
				UserId:              "bar",
			},
			{
				ClientIP:            "8.8.6.6",
				ClientIPCountryCode: "IT",
				UserId:              "foo",
			},
		})
		assert.Nil(t, err)

		_, err = w.Write([]byte(data))
	}))

	cli := mocks.AgentsExternalV1SrvClient{}

	cli.On("SendData", mock.Anything, mock.Anything).Return(&agentsExternalApiV1.SendDataRes{ReceivedAt: timestamppb.Now()}, nil)

	s := services.NewServiceFromConfiguration(&agentsExternalApiV1.AgentConfiguration{
		AgentConfigurationId: "cf:1234567890",
		Name:                 "test",
		Type:                 "sentinel_one",
		CreatedAt:            timestamppb.Now(),
		UpdatedAt:            timestamppb.Now(),
		Data: &agentsExternalApiV1.AgentConfigurationData{
			Data: &agentsExternalApiV1.AgentConfigurationData_ImpossibleTravel{
				ImpossibleTravel: &agentsExternalApiV1.ImpossibleTravelData{
					Url:        ts.URL,
					Token:      "xxx",
					SkipSsl:    true,
					TimeWindow: durationpb.New(time.Hour * 5),
				},
			},
		},
	}, &cli)

	s.RunService(context.Background())

	cli.AssertNumberOfCalls(t, "SendData", 2)

	dFoo, err := json.Marshal(findings{
		Ips: []string{"8.8.8.8", "8.8.6.6"},
		Logs: []map[string]string{
			{
				"ClientIP":              "8.8.8.8",
				"ClientIP_country_code": "US",
				"UserId":                "foo",
			},
			{
				"ClientIP":              "8.8.6.6",
				"ClientIP_country_code": "IT",
				"UserId":                "foo",
			},
		},
		User:      "foo",
		Countries: []string{"US", "IT"},
	})
	assert.Nil(t, err)
	cli.AssertCalled(t, "SendData", mock.Anything, &agentsExternalApiV1.SendDataReq{
		Data:                 string(dFoo),
		DataType:             agentsExternalApiV1.DataType_IMPOSSIBLE_TRAVEL,
		AgentConfigurationId: "cf:1234567890",
	})
}
