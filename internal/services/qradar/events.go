package qradar

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

var DEFAULT_EVENTS_COLUMNS = `QIDNAME(qid), LOGSOURCENAME(logsourceid), CATEGORYNAME(highlevelcategory), CATEGORYNAME(category), PROTOCOLNAME(protocolid), sourceip, sourceport, destinationip, destinationport, QIDDESCRIPTION(qid), username, PROTOCOLNAME(protocolid), RULENAME("creEventList"), sourcegeographiclocation, sourceMAC, sourcev6, destinationgeographiclocation, destinationv6, LOGSOURCETYPENAME(devicetype), credibility, severity, magnitude, eventcount, eventDirection, postNatDestinationIP, postNatDestinationPort, postNatSourceIP, postNatSourcePort, preNatDestinationPort, preNatSourceIP, preNatSourcePort, UTF8(payload), starttime, devicetime`

func (q *QRadarClient) SafeSearchOffenseEvents(ctx context.Context, offenseId int, offenseStartTime int64) []map[string]any {
	res, err := q.SearchOffenseEvents(ctx, offenseId, offenseStartTime)
	if err != nil {
		return []map[string]any{}
	}

	if len(res.Events) == 0 {
		return []map[string]any{}
	}

	return res.Events
}

func (q *QRadarClient) SearchOffenseEvents(ctx context.Context, offenseId int, offenseStartTime int64) (eventSearchResults, error) {
	searchID, err := q.CreateEventsSearch(ctx, DEFAULT_EVENTS_COLUMNS, 20, offenseId, fmt.Sprintf("%d", offenseStartTime))
	if err != nil {
		return eventSearchResults{}, err
	}

	err = q.waitSearchCompleted(ctx, searchID, 10*time.Second)
	if err != nil {
		return eventSearchResults{}, err
	}

	results, err := q.searchResults(ctx, searchID)
	if err != nil {
		return eventSearchResults{}, err
	}

	return results, nil
}

func (q *QRadarClient) waitSearchCompleted(ctx context.Context, searchID string, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	for {
		searchResponse, err := q.searchStatus(ctx, searchID)
		if err != nil {
			return err
		}

		if searchResponse.Status == qradarSearchStatusCompleted {
			return nil
		}

		if searchResponse.Status == qradarSearchStatusCanceled {
			return fmt.Errorf("search canceled")
		}

		if ctx.Err() != nil {
			return ctx.Err()
		}

		time.Sleep(1 * time.Second)
	}
}

// CreateEventsSearch creates a search for events related to a specific offense
func (q *QRadarClient) CreateEventsSearch(
	ctx context.Context,
	// fetchMode FetchMode,
	eventsColumns string,
	eventsLimit int,
	offenseID int,
	offenseStartTime string,
) (string, error) {
	queryExpression := fmt.Sprintf(
		"SELECT %s FROM events WHERE INOFFENSE(%d) limit %d START %s",
		eventsColumns,
		offenseID,
		eventsLimit,
		offenseStartTime,
	)

	searchResponse, err := q.searchCreate(ctx, queryExpression)
	if err != nil {
		return "", fmt.Errorf("search creation failed: %w", err)
	}

	if searchResponse.SearchID == "" {
		return "", fmt.Errorf("empty search ID in response")
	}

	return searchResponse.SearchID, nil
}

type SearchResponse struct {
	SearchID string `json:"search_id"`
}

// searchCreate creates a new search in QRadar
func (q *QRadarClient) searchCreate(ctx context.Context, queryExpression string) (*SearchResponse, error) {
	path, err := url.JoinPath(q.url, "/api/ariel/searches")
	if err != nil {
		return nil, err
	}

	params := url.Values{}
	params.Add("query_expression", queryExpression)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, path, nil)
	if err != nil {
		return nil, err
	}

	qp := req.URL.Query()

	qp.Add("query_expression", queryExpression)
	req.URL.RawQuery = qp.Encode()

	q.addBaseHeadersToRequest(req)
	req.Header.Set("Content-Type", "application/json")

	res, err := q.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	payload, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("status code %d", res.StatusCode)
	}

	var searchResponse SearchResponse
	err = json.Unmarshal(payload, &searchResponse)
	if err != nil {
		return nil, err
	}

	return &searchResponse, nil
}

type qradarSearchStatus string

const (
	qradarSearchStatusCanceled  qradarSearchStatus = "CANCELED"
	qradarSearchStatusCompleted qradarSearchStatus = "COMPLETED"
)

type searchStatusResponse struct {
	SearchID string             `json:"search_id"`
	Status   qradarSearchStatus `json:"status"`
}

// searchStatus retrieves the status of a search
func (q *QRadarClient) searchStatus(ctx context.Context, searchID string) (*searchStatusResponse, error) {
	path, err := url.JoinPath(q.url, "/api/ariel/searches", searchID)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, path, nil)
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

	var searchStatusResponse searchStatusResponse
	err = json.Unmarshal(payload, &searchStatusResponse)
	if err != nil {
		return nil, err
	}

	return &searchStatusResponse, nil
}

type eventSearchResults struct {
	Events []map[string]any `json:"events"`
}

func (q *QRadarClient) searchResults(ctx context.Context, searchID string) (eventSearchResults, error) {
	path, err := url.JoinPath(q.url, "/api/ariel/searches", searchID, "results")
	if err != nil {
		return eventSearchResults{}, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, path, nil)
	if err != nil {
		return eventSearchResults{}, err
	}

	q.addBaseHeadersToRequest(req)

	res, err := q.client.Do(req)
	if err != nil {
		return eventSearchResults{}, err
	}

	defer res.Body.Close()

	payload, err := io.ReadAll(res.Body)
	if err != nil {
		return eventSearchResults{}, err
	}

	var results eventSearchResults
	err = json.Unmarshal(payload, &results)
	if err != nil {
		return eventSearchResults{}, err
	}

	return results, nil
}
