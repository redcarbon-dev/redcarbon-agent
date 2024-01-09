package qradar

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type QRadarClient struct {
	client *http.Client
	token  string
	url    string
}

func NewQRadarClient(url, token string, skipSSL bool) QRadarClient {
	return QRadarClient{
		token: token,
		url:   url,
		client: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: !skipSSL},
			},
		},
	}
}

func (q *QRadarClient) addBaseHeadersToRequest(req *http.Request) {
	req.Header.Set("Accept", "application/json")
	req.Header.Set("SEC", q.token)
	req.Header.Set("X-Requested-By", "RedCarbon Agent")
}

func (q *QRadarClient) FetchOffenses(ctx context.Context, start, end time.Time) ([]map[string]interface{}, error) {
	rFetchOffenses, err := url.JoinPath(q.url, "/api/siem/offenses")
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rFetchOffenses, nil)
	if err != nil {
		return nil, err
	}

	q.addBaseHeadersToRequest(req)
	addTimeFilteringToRequest(req, start, end)

	res, err := q.client.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code %d", res.StatusCode)
	}

	defer res.Body.Close()

	payload, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var offenses []map[string]interface{}

	err = json.Unmarshal(payload, &offenses)
	if err != nil {
		return nil, err
	}

	return offenses, nil
}

func (q *QRadarClient) FetchSourceAddresses(ctx context.Context, sourceAddressIds []int) ([]map[string]interface{}, error) {
	rFetchSourceAd, err := url.JoinPath(q.url, "/api/siem/source_addresses")
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rFetchSourceAd, nil)
	if err != nil {
		return nil, err
	}

	addIdsFilteringToRequest(req, sourceAddressIds)
	q.addBaseHeadersToRequest(req)

	res, err := q.client.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code %d", res.StatusCode)
	}

	defer res.Body.Close()

	payload, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var sourceAddresses []map[string]interface{}

	err = json.Unmarshal(payload, &sourceAddresses)
	if err != nil {
		return nil, err
	}

	return sourceAddresses, nil
}

func (q *QRadarClient) FetchLocalDestinationAddresses(ctx context.Context, localDestinationAddressIds []int) ([]map[string]interface{}, error) {
	rFetchLocalDestAd, err := url.JoinPath(q.url, "/api/siem/local_destination_addresses")
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rFetchLocalDestAd, nil)
	if err != nil {
		return nil, err
	}

	addIdsFilteringToRequest(req, localDestinationAddressIds)
	q.addBaseHeadersToRequest(req)

	res, err := q.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code %d", res.StatusCode)
	}

	payload, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var localDestinationAddresses []map[string]interface{}

	err = json.Unmarshal(payload, &localDestinationAddresses)
	if err != nil {
		return nil, err
	}

	return localDestinationAddresses, nil
}

func (q *QRadarClient) FetchOffenseType(ctx context.Context, offenseTypeId int) (map[string]interface{}, error) {
	rFetchOffenseType, err := url.JoinPath(q.url, fmt.Sprintf("/api/siem/offense_types/%d", offenseTypeId))
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rFetchOffenseType, nil)
	if err != nil {
		return nil, err
	}

	q.addBaseHeadersToRequest(req)

	res, err := q.client.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code %d", res.StatusCode)
	}

	defer res.Body.Close()

	payload, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var offenseType map[string]interface{}

	err = json.Unmarshal(payload, &offenseType)
	if err != nil {
		return nil, err
	}

	return offenseType, nil
}

func (q *QRadarClient) RetrieveOffenseUrl(offenseId int) *string {
	cUrl, err := url.JoinPath(q.url, fmt.Sprintf("/console/do/sem/offensesummary?appName=Sem&pageId=OffenseSummary&summaryId=%d", offenseId))
	if err != nil {
		return nil
	}

	return &cUrl
}

func addIdsFilteringToRequest(req *http.Request, ids []int) {
	var idsStr []string

	for _, id := range ids {
		idsStr = append(idsStr, fmt.Sprintf("id=%d", id))
	}

	filter := strings.Join(idsStr, " OR ")

	q := req.URL.Query()

	q.Add("filter", filter)

	req.URL.RawQuery = q.Encode()

	return
}

func addTimeFilteringToRequest(req *http.Request, start, end time.Time) {
	q := req.URL.Query()

	q.Add("filter", fmt.Sprintf("start_time>=%d AND start_time<=%d", start.UnixMilli(), end.UnixMilli()))
	q.Add("sort", "-start_time")

	req.URL.RawQuery = q.Encode()

	return
}
