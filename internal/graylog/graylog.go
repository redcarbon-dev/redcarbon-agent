package graylog

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/gocarina/gocsv"
)

type Client struct {
	url     string
	token   string
	skipSSL bool
}

func NewGrayLogClient(token string, url string, skipSSL bool) Client {
	return Client{
		token:   base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:token", token))),
		url:     url,
		skipSSL: skipSSL,
	}
}

type queryString struct {
	Type        string `json:"type"`
	QueryString string `json:"query_string"`
}

type timeRange struct {
	Type string    `json:"type"`
	From time.Time `json:"from"`
	To   time.Time `json:"to"`
}

type payload struct {
	QueryString   queryString `json:"query_string"`
	TimeRange     timeRange   `json:"timerange"`
	FieldsInOrder []string    `json:"fields_in_order"`
}

func (c Client) QueryData(ctx context.Context, query string, from time.Time, to time.Time, fields []string) ([]map[string]string, error) {
	p := payload{
		QueryString: queryString{
			QueryString: query,
			Type:        "elasticsearch",
		},
		TimeRange: timeRange{
			Type: "absolute",
			From: from,
			To:   to,
		},
		FieldsInOrder: fields,
	}

	jP, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}

	rUrl, err := url.JoinPath(c.url, "/api/views/search/messages")
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, rUrl, bytes.NewReader(jP))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Basic %s", c.token))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Requested-By", "RC Agent")

	client := http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: c.skipSSL}}}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	pRes, err := gocsv.CSVToMaps(res.Body)
	if err != nil {
		return nil, err
	}

	return pRes, nil
}
