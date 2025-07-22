package routines

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"time"

	"connectrpc.com/connect"
	"github.com/sirupsen/logrus"
	agents_publicv1 "pkg.redcarbon.ai/proto/redcarbon/agents_public/v1"
)

var httpCli = http.Client{
	Timeout: 10 * time.Second,
}

func (r RoutineConfig) ProxyRoutine(ctx context.Context) {
	logrus.Info("Starting the proxy routine...")

	// Prepare the request to fetch agent requests
	req := connect.NewRequest(&agents_publicv1.FetchAgentRequestsRequest{})
	req.Header().Set("authorization", fmt.Sprintf("ApiToken %s", r.profile.Profile.Token))

	ctxWithTimeout, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	logrus.Info("Fetching the agent requests...")
	res, err := r.agentsCli.FetchAgentRequests(ctxWithTimeout, req)
	if err != nil {
		logrus.WithError(err).Error("Error while fetching the agent requests")
		return
	}

	logrus.Infof("Running %d agent requests...", len(res.Msg.Requests))

	var wg sync.WaitGroup
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt)
	defer stop()

	for _, agentReq := range res.Msg.Requests {
		r.processRequest(ctx, agentReq, &wg)
	}

	// Wait for all requests to finish or context to be cancelled
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-ctx.Done():
		logrus.Warn("Proxy routine interrupted by signal or context cancellation")
	case <-done:
		logrus.Info("Proxy routine completed")
	}
}

func (r RoutineConfig) processRequest(ctx context.Context, req *agents_publicv1.AgentRequest, wg *sync.WaitGroup) {
	wg.Add(1)

	go func() {
		defer wg.Done()
		l := logrus.WithFields(logrus.Fields{
			"requestId": req.RequestId,
		})

		l.Info("Handling request...")

		httpReq, err := r.createHTTPProxyRequest(ctx, req)
		if err != nil {
			l.WithError(err).Error("Error while creating the HTTP request")
			return
		}

		l.Infof("Executing HTTP request: %s %s", req.Method, req.Url)

		httpRes, err := httpCli.Do(httpReq)
		if err != nil {
			l.WithError(err).Error("Error while executing the HTTP request")
			return
		}
		defer httpRes.Body.Close()

		l.Infof("Request completed with status code %d, sending the response to the server", httpRes.StatusCode)

		err = r.sendResponseToServer(ctx, req, httpRes)
		if err != nil {
			l.WithError(err).Error("Error while sending the response to the server")
			return
		}

		l.Info("Response sent to the server")
	}()
}

func (r RoutineConfig) createHTTPProxyRequest(ctx context.Context, req *agents_publicv1.AgentRequest) (*http.Request, error) {
	// Url is already validated by the server
	httpReq, err := http.NewRequestWithContext(ctx, req.Method, req.Url, bytes.NewReader(req.Body))
	if err != nil {
		return nil, err
	}

	for key, value := range req.Headers {
		httpReq.Header.Set(key, value)
	}

	return httpReq, nil
}

func (r RoutineConfig) sendResponseToServer(ctx context.Context, req *agents_publicv1.AgentRequest, httpRes *http.Response) error {
	body, err := io.ReadAll(httpRes.Body)
	if err != nil {
		return err
	}

	headers := make(map[string]string)
	for key, values := range httpRes.Header {
		headers[key] = strings.Join(values, ",")
	}

	response := connect.NewRequest(&agents_publicv1.SubmitAgentResponseRequest{
		RequestId: req.RequestId,
		Response: &agents_publicv1.AgentResponse{
			Status:  int32(httpRes.StatusCode),
			Body:    body,
			Headers: headers,
		},
	})
	response.Header().Set("authorization", fmt.Sprintf("ApiToken %s", r.profile.Profile.Token))

	_, err = r.agentsCli.SubmitAgentResponse(ctx, response)
	if err != nil {
		return err
	}

	return nil
}
