// Package github is a minimal client for the GitHub REST API endpoints the agent uses.
package github

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"pkg.redcarbon.ai/internal/build"
)

const (
	DefaultBaseURL = "https://api.github.com"
	apiVersion     = "2022-11-28"
	defaultTimeout = 30 * time.Second
)

// Client is a minimal client for the public, unauthenticated GitHub REST API.
type Client struct {
	httpClient *http.Client
	baseURL    string
}

// Release is a GitHub release with its downloadable assets.
type Release struct {
	TagName string  `json:"tag_name"`
	Assets  []Asset `json:"assets"`
}

// Asset is a file attached to a GitHub release.
type Asset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

// APIError is returned when GitHub answers with a non-2xx status.
type APIError struct {
	StatusCode int
	Message    string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("github API returned %d: %s", e.StatusCode, e.Message)
}

// NewClient returns a Client for baseURL; a nil httpClient gets a default one with a timeout.
func NewClient(httpClient *http.Client, baseURL string) *Client {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: defaultTimeout}
	}

	return &Client{httpClient: httpClient, baseURL: baseURL}
}

// LatestRelease returns the latest published, non-prerelease release of owner/repo.
func (c *Client) LatestRelease(ctx context.Context, owner, repo string) (*Release, error) {
	u, err := url.JoinPath(c.baseURL, "repos", owner, repo, "releases", "latest")
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", apiVersion)
	// GitHub rejects API requests without a User-Agent.
	req.Header.Set("User-Agent", build.UserAgent())

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, newAPIError(resp)
	}

	var rel Release
	if err := json.NewDecoder(resp.Body).Decode(&rel); err != nil {
		return nil, fmt.Errorf("decoding latest release: %w", err)
	}

	return &rel, nil
}

func newAPIError(resp *http.Response) *APIError {
	apiErr := &APIError{StatusCode: resp.StatusCode, Message: resp.Status}

	var body struct {
		Message string `json:"message"`
	}
	if err := json.NewDecoder(io.LimitReader(resp.Body, 64<<10)).Decode(&body); err == nil && body.Message != "" {
		apiErr.Message = body.Message
	}

	return apiErr
}
