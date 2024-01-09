package sentinel_one

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type SentinelOneClient struct {
	client *http.Client
	token  string
	url    string
}

var threatsLimit = 100

func NewSentinelOneClient(url, token string, skipSSL bool) SentinelOneClient {
	return SentinelOneClient{
		token: token,
		url:   url,
		client: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: !skipSSL},
			},
		},
	}
}

func (s *SentinelOneClient) addBaseHeadersToRequest(req *http.Request) {
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("ApiToken %s", s.token))
	req.Header.Set("X-Requested-By", "RedCarbon Agent")
}

func (s *SentinelOneClient) FetchThreats(ctx context.Context, start, end time.Time) ([]map[string]interface{}, error) {
	page := 0
	threats := []map[string]interface{}{}

	rFetchThreats, err := url.JoinPath(s.url, "/web/api/v2.1/threats")
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rFetchThreats, nil)
	if err != nil {
		return nil, err
	}

	s.addBaseHeadersToRequest(req)
	q := req.URL.Query()
	q.Add("createdAt__gt", start.UTC().Format(time.RFC3339))
	q.Add("createdAt__lt", end.UTC().Format(time.RFC3339))

	for {
		q.Set("skip", fmt.Sprintf("%d", page))
		q.Set("limit", fmt.Sprintf("%d", threatsLimit))
		req.URL.RawQuery = q.Encode()

		res, err := s.client.Do(req)
		if err != nil {
			return nil, err
		}

		response, err := parseResponse(res)
		if err != nil {
			return nil, err
		}

		if len(response.Data) == 0 {
			return threats, nil
		}

		threats = append(threats, response.Data...)
		page += threatsLimit
	}
}

func parseResponse(res *http.Response) (SentinelOneFetchThreatsResponse, error) {
	var response SentinelOneFetchThreatsResponse

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return SentinelOneFetchThreatsResponse{}, fmt.Errorf("status code %d", res.StatusCode)
	}

	err := json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return SentinelOneFetchThreatsResponse{}, err
	}

	return response, nil
}
