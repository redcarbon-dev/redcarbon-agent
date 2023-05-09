package graylog_datamine_test

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gocarina/gocsv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/protobuf/types/known/timestamppb"

	"pkg.redcarbon.ai/internal/services"
	"pkg.redcarbon.ai/mocks"
	agentsPublicApiV1 "pkg.redcarbon.ai/proto/redcarbon/public_apis/agents/api/v1"
)

type grayLogDataMine struct {
	Timestamp string `csv:"timestamp"`
	Source    string `csv:"source"`
	Uuid      string `csv:"gl2_message_id"`
	Message   string `csv:"message"`
}

func TestNewGrayLogDataMineService(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		dToken, err := base64.StdEncoding.DecodeString(strings.Replace(r.Header.Get("Authorization"), "Basic ", "", -1))
		assert.Nil(t, err)
		assert.Equal(t, "xxx:token", string(dToken))

		data, err := gocsv.MarshalString([]grayLogDataMine{
			{
				Timestamp: "2021-03-23T09:04:49.000Z",
				Source:    "127.0.0.1",
				Uuid:      "18515d34-112d-44c2-ad98-8fb37d387dbd",
				Message:   "Test test",
			},
			{
				Timestamp: "2022-03-23T09:04:49.000Z",
				Source:    "127.0.0.1",
				Uuid:      "bc5d1607-a7ec-46df-8d98-5ebfe5b1e7d6",
				Message:   "Test test 2",
			},
		})
		assert.Nil(t, err)

		_, err = w.Write([]byte(data))
		assert.Nil(t, err)
	}))

	cli := mocks.AgentsPublicApiV1SrvClient{}

	cli.On("GetGrayLogDataMinePendingQueries", mock.Anything, mock.Anything).Return(&agentsPublicApiV1.GetGrayLogDataMinePendingQueriesRes{
		GraylogDatamineQueries: []*agentsPublicApiV1.GrayLogDataMineQuery{
			{
				Id:              "0",
				SearchStartTime: timestamppb.New(time.Now().Add(-time.Hour)),
				SearchStopTime:  timestamppb.Now(),
				Query:           "",
			},
		},
	}, nil)
	cli.On("SendGrayLogDatamineQueryResultsData", mock.Anything, mock.Anything).Return(&agentsPublicApiV1.SendGrayLogDatamineQueryResultsDataRes{ReceivedAt: timestamppb.Now()}, nil)

	s := services.NewServiceFromConfiguration(&agentsPublicApiV1.AgentConfiguration{
		AgentConfigurationId: "cf:1234567890",
		Name:                 "test",
		Type:                 "sentinel_one",
		CreatedAt:            timestamppb.Now(),
		UpdatedAt:            timestamppb.Now(),
		Data: &agentsPublicApiV1.AgentConfigurationData{
			Data: &agentsPublicApiV1.AgentConfigurationData_GraylogDatamine{
				GraylogDatamine: &agentsPublicApiV1.GrayLogDataMineData{
					Url:     ts.URL,
					Token:   "xxx",
					SkipSsl: true,
				},
			},
		},
	}, &cli)

	s.RunService(context.Background())

	cli.AssertNumberOfCalls(t, "SendGrayLogDatamineQueryResultsData", 1)
	cli.AssertNotCalled(t, "SendGrayLogDatamineQueryErrorData", mock.Anything, mock.Anything)
}

func TestShouldSendTheErrorInCaseSomethingWentWrong(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		dToken, err := base64.StdEncoding.DecodeString(strings.Replace(r.Header.Get("Authorization"), "Basic ", "", -1))
		assert.Nil(t, err)
		assert.Equal(t, "xxx:token", string(dToken))

		data, err := gocsv.MarshalString([]grayLogDataMine{
			{
				Timestamp: "2021-03-23T09:04:49.000Z",
				Source:    "127.0.0.1",
				Uuid:      "18515d34-112d-44c2-ad98-8fb37d387dbd",
				Message:   "Test test",
			},
			{
				Timestamp: "2022-03-23T09:04:49.000Z",
				Source:    "127.0.0.1",
				Uuid:      "bc5d1607-a7ec-46df-8d98-5ebfe5b1e7d6",
				Message:   "Test test 2",
			},
		})
		assert.Nil(t, err)

		_, err = w.Write([]byte(data))
		assert.Nil(t, err)
	}))

	cli := mocks.AgentsPublicApiV1SrvClient{}

	cli.On("GetGrayLogDataMinePendingQueries", mock.Anything, mock.Anything).Return(&agentsPublicApiV1.GetGrayLogDataMinePendingQueriesRes{
		GraylogDatamineQueries: []*agentsPublicApiV1.GrayLogDataMineQuery{
			{
				Id:              "0",
				SearchStartTime: timestamppb.New(time.Now().Add(-time.Hour)),
				SearchStopTime:  timestamppb.Now(),
				Query:           "",
			},
			{
				Id:              "1",
				SearchStartTime: timestamppb.New(time.Now().Add(-time.Hour)),
				SearchStopTime:  timestamppb.Now(),
				Query:           "",
			},
		},
	}, nil)
	cli.On("SendGrayLogDatamineQueryResultsData", mock.Anything, mock.Anything).Return(&agentsPublicApiV1.SendGrayLogDatamineQueryResultsDataRes{ReceivedAt: timestamppb.Now()}, nil).Once()
	cli.On("SendGrayLogDatamineQueryResultsData", mock.Anything, mock.Anything).Return(nil, fmt.Errorf("test error")).Once()
	cli.On("SendGrayLogDatamineQueryErrorData", mock.Anything, mock.Anything).Return(&agentsPublicApiV1.SendGrayLogDatamineQueryErrorDataRes{ReceivedAt: timestamppb.Now()}, nil)

	s := services.NewServiceFromConfiguration(&agentsPublicApiV1.AgentConfiguration{
		AgentConfigurationId: "cf:1234567890",
		Name:                 "test",
		Type:                 "sentinel_one",
		CreatedAt:            timestamppb.Now(),
		UpdatedAt:            timestamppb.Now(),
		Data: &agentsPublicApiV1.AgentConfigurationData{
			Data: &agentsPublicApiV1.AgentConfigurationData_GraylogDatamine{
				GraylogDatamine: &agentsPublicApiV1.GrayLogDataMineData{
					Url:     ts.URL,
					Token:   "xxx",
					SkipSsl: true,
				},
			},
		},
	}, &cli)

	s.RunService(context.Background())

	cli.AssertNumberOfCalls(t, "SendGrayLogDatamineQueryResultsData", 2)
	cli.AssertNumberOfCalls(t, "SendGrayLogDatamineQueryErrorData", 1)
}
