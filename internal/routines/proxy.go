package routines

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"connectrpc.com/connect"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/transform"
	agents_publicv1 "pkg.redcarbon.ai/proto/redcarbon/agents_public/v1"
)

var defaultTimeout = 1 * time.Minute

var httpCli = http.Client{
	Timeout: defaultTimeout,
	Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	},
}

func (r RoutineConfig) ProxyRoutine(ctx context.Context) {
	logrus.Info("Starting the proxy routine...")

	// Prepare the request to fetch agent requests
	req := connect.NewRequest(&agents_publicv1.FetchAgentRequestsRequest{})
	req.Header().Set("authorization", fmt.Sprintf("ApiToken %s", r.profile.Profile.Token))

	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	logrus.Info("Fetching the agent requests...")
	res, err := r.agentsCli.FetchAgentRequests(ctx, req)
	if err != nil {
		logrus.WithError(err).Error("Error while fetching the agent requests")
		return
	}

	logrus.Infof("Running %d agent requests...", len(res.Msg.Requests))

	var wg sync.WaitGroup

	for _, agentReq := range res.Msg.Requests {
		wg.Add(1)
		go func() {
			defer wg.Done()
			r.processRequest(ctx, agentReq)
		}()
	}

	wg.Wait()

	logrus.Info("Proxy routine completed")
}

func (r RoutineConfig) processRequest(ctx context.Context, req *agents_publicv1.AgentRequest) {
	l := logrus.WithFields(logrus.Fields{
		"requestId": req.RequestId,
	})

	l.Info("Handling request...")

	httpReq, err := r.createHTTPProxyRequest(ctx, req)
	if err != nil {
		l.WithError(err).Error("Error while creating the HTTP request")
		r.sendErrorToServer(ctx, req, "error while creating the HTTP request")
		return
	}

	l.Infof("Executing HTTP request: %s %s", req.Method, req.Url)

	httpRes, err := httpCli.Do(httpReq)
	if err != nil {
		l.WithError(err).Error("Error while executing the HTTP request")
		r.sendErrorToServer(ctx, req, "error while executing the HTTP request")
		return
	}
	defer httpRes.Body.Close()

	l.Infof("Request completed with status code %d, sending the response to the server", httpRes.StatusCode)

	err = r.sendResponseToServer(ctx, req, httpRes)
	if err != nil {
		l.WithError(err).Error("Error while sending the response to the server")
		r.sendErrorToServer(ctx, req, "error while sending the response to the server")
		return
	}
}

func (r RoutineConfig) createHTTPProxyRequest(ctx context.Context, req *agents_publicv1.AgentRequest) (*http.Request, error) {
	// Url is already validated by the server
	httpReq, err := http.NewRequestWithContext(ctx, req.Method, req.Url, bytes.NewBuffer(req.Body))
	if err != nil {
		return nil, err
	}

	for key, value := range req.Headers {
		httpReq.Header.Set(key, value)
	}

	return httpReq, nil
}

func (r RoutineConfig) sendErrorToServer(ctx context.Context, req *agents_publicv1.AgentRequest, reason string) {
	body, err := json.Marshal(map[string]string{
		"error": reason,
	})
	if err != nil {
		logrus.WithError(err).Error("Error while marshalling the error to the server")
		return
	}

	response := connect.NewRequest(&agents_publicv1.SubmitAgentResponseRequest{
		RequestId: req.RequestId,
		Response: &agents_publicv1.AgentResponse{
			Status:  int32(http.StatusInternalServerError),
			Body:    body,
			Headers: make(map[string]string),
		},
	})
	response.Header().Set("authorization", fmt.Sprintf("ApiToken %s", r.profile.Profile.Token))

	_, err = r.agentsCli.SubmitAgentResponse(ctx, response)
	if err != nil {
		logrus.WithError(err).Error("Error while sending the error to the server")
	}
}

func (r RoutineConfig) sendResponseToServer(ctx context.Context, req *agents_publicv1.AgentRequest, httpRes *http.Response) error {
	body, err := io.ReadAll(httpRes.Body)
	if err != nil {
		return err
	}

	// Handle charset conversion
	convertedBody, err := r.convertCharset(body, httpRes.Header.Get("Content-Type"))
	if err != nil {
		logrus.WithError(err).Warn("Failed to convert charset, using original body")
		convertedBody = body
	}

	logrus.Infof("Response body: %s", string(convertedBody))

	headers := make(map[string]string)
	for key, values := range httpRes.Header {
		headers[key] = strings.Join(values, ",")
	}

	response := connect.NewRequest(&agents_publicv1.SubmitAgentResponseRequest{
		RequestId: req.RequestId,
		Response: &agents_publicv1.AgentResponse{
			Status:  int32(httpRes.StatusCode),
			Body:    convertedBody,
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

// convertCharset converts the response body to UTF-8 based on the Content-Type charset
func (r RoutineConfig) convertCharset(body []byte, contentType string) ([]byte, error) {
	if len(body) == 0 {
		return body, nil
	}

	// Extract charset from Content-Type header
	_, name, _ := charset.DetermineEncoding(body, contentType)

	// If no charset is detected or it's already UTF-8, return as-is
	if name == "utf-8" || name == "" {
		return body, nil
	}

	// Get the encoding for the detected charset
	e, _ := charset.Lookup(name)
	if e == nil {
		// If we can't find the encoding, return the original body
		return body, nil
	}

	// Convert to UTF-8
	reader := transform.NewReader(bytes.NewReader(body), e.NewDecoder())
	converted, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to convert charset %s to UTF-8: %w", name, err)
	}

	return converted, nil
}
