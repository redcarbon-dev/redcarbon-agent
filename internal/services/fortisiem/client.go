package fortisiem

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type FortiSIEMClient struct {
	client   *http.Client
	username string
	password string
	url      string
}

func NewFortiSIEMClient(url, username, password string, verifySSL bool) FortiSIEMClient {
	return FortiSIEMClient{
		username: username,
		password: password,
		url:      url,
		client: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: !verifySSL},
			},
		},
	}
}

func (f *FortiSIEMClient) FetchAlerts(ctx context.Context, start, end time.Time) ([]map[string]interface{}, error) {
	req, err := f.createHTTPRequest(ctx, http.MethodGet, "/phoenix/rest/pub/incident")
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()

	q.Add("timeFrom", strconv.FormatInt(start.UTC().UnixMilli(), 10))
	q.Add("timeTo", strconv.FormatInt(end.UTC().UnixMilli(), 10))

	req.URL.RawQuery = q.Encode()

	res, err := f.client.Do(req)
	if err != nil {
		return nil, err
	}

	response, err := parseResponse(res)
	if err != nil {
		return nil, err
	}

	if response.Pages == 1 {
		return response.Data, nil
	}

	alerts := response.Data
	queryId := response.QueryID

	for page := 2; page <= response.Pages; page++ {
		reqP, err := f.createHTTPRequest(ctx, http.MethodGet, fmt.Sprintf("/phoenix/rest/pub/incident/%s/%d", queryId, page))
		if err != nil {
			return nil, err
		}

		resP, err := f.client.Do(reqP)
		if err != nil {
			return nil, err
		}

		responseP, err := parseResponse(resP)
		if err != nil {
			return nil, err
		}

		alerts = append(alerts, responseP.Data...)
	}

	return alerts, nil
}

func parseResponse(res *http.Response) (FortiSIEMFetchAlertsResponse, error) {
	var response FortiSIEMFetchAlertsResponse

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return FortiSIEMFetchAlertsResponse{}, fmt.Errorf("status code %d", res.StatusCode)
	}

	err := json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return FortiSIEMFetchAlertsResponse{}, err
	}

	return response, nil
}

func (f *FortiSIEMClient) createHTTPRequest(ctx context.Context, method string, path string) (*http.Request, error) {
	r, err := url.JoinPath(f.url, path)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, method, r, nil)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(f.username, f.password)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-Requested-By", "RedCarbon Agent")

	return req, nil
}
