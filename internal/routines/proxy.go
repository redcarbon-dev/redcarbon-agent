package routines

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"sync"
	"time"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/durationpb"
	agents_publicv1 "pkg.redcarbon.ai/proto/redcarbon/agents_public/v1"
)

var defaultTimeout = 1 * time.Minute

var httpCli = &http.Client{
	Timeout: defaultTimeout,
	Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	},
}

func (r RoutineConfig) ProxyRoutine(ctx context.Context) {
	r.Logger().Debug("Starting the proxy routine...")

	// Prepare the request to fetch agent requests
	req := connect.NewRequest(&agents_publicv1.FetchAgentRequestsRequest{})
	req.Header().Set("authorization", fmt.Sprintf("ApiToken %s", r.profile.Profile.Token))

	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	r.Logger().Debug("Fetching the agent requests...")
	res, err := r.agentsCli.FetchAgentRequests(ctx, req)
	if err != nil {
		r.Logger().Error("Error while fetching the agent requests", "error", err)
		return
	}

	r.Logger().Debug("Running agent requests...", "count", len(res.Msg.Requests))

	var wg sync.WaitGroup

	for _, agentReq := range res.Msg.Requests {
		wg.Add(1)
		go func() {
			defer wg.Done()
			r.processRequest(ctx, agentReq)
		}()
	}

	wg.Wait()

	r.Logger().Debug("Proxy routine completed")
}

func (r RoutineConfig) processRequest(ctx context.Context, req *agents_publicv1.AgentRequest) {
	l := r.Logger().With("id", req.RequestId)

	l.Debug("Handling request...")

	ctx, cancel := extractTimeout(ctx, req.Timeout)
	defer cancel()

	httpReq, err := r.createHTTPProxyRequest(ctx, req)
	if err != nil {
		l.Error("Error while creating the HTTP request", "error", err)
		r.sendErrorToServer(ctx, req, "error while creating the HTTP request")
		return
	}

	l.Info("Executing HTTP request",
		"method", req.Method,
		"url", redactedURL(req.Url),
		"headers", headerNames(httpReq.Header),
		"timeout", req.Timeout.String(),
	)

	httpRes, err := httpCli.Do(httpReq)
	if err != nil {
		l.Error("Error while executing the HTTP request", "error", err)
		r.sendErrorToServer(ctx, req, "error while executing the HTTP request")
		return
	}
	defer httpRes.Body.Close()

	l.Info("Request completed, sending the response to the server", "status", httpRes.StatusCode)

	err = r.sendResponseToServer(ctx, req, httpRes)
	if err != nil {
		l.Error("Error while sending the response to the server", "error", err)
		r.sendErrorToServer(ctx, req, "error while sending the response to the server")
		return
	}
}

func extractTimeout(ctx context.Context, requestTimeout *durationpb.Duration) (context.Context, context.CancelFunc) {
	timeout := 10 * time.Second
	if requestTimeout != nil && requestTimeout.AsDuration() > 0 {
		timeout = requestTimeout.AsDuration()
	}

	return context.WithTimeout(ctx, timeout)
}

func (r RoutineConfig) createHTTPProxyRequest(ctx context.Context, req *agents_publicv1.AgentRequest) (*http.Request, error) {
	// Url is already validated by the server
	httpReq, err := http.NewRequestWithContext(ctx, req.Method, req.Url, bytes.NewBuffer(req.Body))
	if err != nil {
		return nil, err
	}

	for key, value := range req.Headers {
		if strings.ToLower(key) == "accept-encoding" {
			continue
		}

		for _, v := range value.Values {
			httpReq.Header.Add(key, v)
		}
	}

	return httpReq, nil
}

// headerNames returns the sorted header names only: values carry vendor credentials
// (Authorization, API keys, cookies) and must never reach the agent log.
func headerNames(h http.Header) []string {
	names := make([]string, 0, len(h))
	for name := range h {
		names = append(names, name)
	}
	slices.Sort(names)
	return names
}

// redactedURL masks the userinfo password and every query value, keeping the
// parameter names: some vendors take API keys in the query string.
func redactedURL(raw string) string {
	u, err := url.Parse(raw)
	if err != nil {
		return "<unparseable URL>"
	}
	q := u.Query()
	for name := range q {
		q[name] = []string{"xxxxx"}
	}
	u.RawQuery = q.Encode()
	return u.Redacted()
}

func (r RoutineConfig) sendErrorToServer(ctx context.Context, req *agents_publicv1.AgentRequest, reason string) {
	body, err := json.Marshal(map[string]string{
		"error": reason,
	})
	if err != nil {
		r.Logger().Error("Error while marshalling the error to the server", "error", err)
		return
	}

	response := connect.NewRequest(&agents_publicv1.SubmitAgentResponseRequest{
		RequestId: req.RequestId,
		Response: &agents_publicv1.AgentResponse{
			Status:  int32(http.StatusInternalServerError),
			Body:    body,
			Headers: map[string]*agents_publicv1.ValueHeader{},
		},
	})
	response.Header().Set("authorization", fmt.Sprintf("ApiToken %s", r.profile.Profile.Token))

	_, err = r.agentsCli.SubmitAgentResponse(ctx, response)
	if err != nil {
		r.Logger().Error("Error while sending the error to the server", "error", err)
	}
}

func (r RoutineConfig) sendResponseToServer(ctx context.Context, req *agents_publicv1.AgentRequest, httpRes *http.Response) error {
	body, err := io.ReadAll(httpRes.Body)
	if err != nil {
		return err
	}

	headers := make(map[string]*agents_publicv1.ValueHeader)
	for key, values := range httpRes.Header {
		headers[key] = &agents_publicv1.ValueHeader{
			Values: values,
		}
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
