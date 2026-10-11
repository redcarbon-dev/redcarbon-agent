package github

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"pkg.redcarbon.ai/internal/build"
)

func TestLatestRelease(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/repos/redcarbon-dev/redcarbon-agent/releases/latest", r.URL.Path)
		assert.Equal(t, "application/vnd.github+json", r.Header.Get("Accept"))
		assert.Equal(t, "2022-11-28", r.Header.Get("X-GitHub-Api-Version"))
		assert.Equal(t, "redcarbon-agent/"+build.Version, r.Header.Get("User-Agent"))

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"tag_name": "v1.2.3",
			"name": "v1.2.3",
			"assets": [
				{"name": "redcarbon_1.2.3_linux_amd64.tar.gz", "browser_download_url": "https://example.com/a.tar.gz", "size": 42},
				{"name": "checksums.txt", "browser_download_url": "https://example.com/checksums.txt"}
			]
		}`))
	}))
	defer srv.Close()

	rel, err := NewClient(srv.Client(), srv.URL).LatestRelease(t.Context(), "redcarbon-dev", "redcarbon-agent")
	require.NoError(t, err)

	assert.Equal(t, "v1.2.3", rel.TagName)
	assert.Equal(t, []Asset{
		{Name: "redcarbon_1.2.3_linux_amd64.tar.gz", BrowserDownloadURL: "https://example.com/a.tar.gz"},
		{Name: "checksums.txt", BrowserDownloadURL: "https://example.com/checksums.txt"},
	}, rel.Assets)
}

func TestLatestReleaseNon2xx(t *testing.T) {
	tests := []struct {
		name    string
		status  int
		body    string
		wantMsg string
	}{
		{"no release", http.StatusNotFound, `{"message":"Not Found"}`, "Not Found"},
		{"rate limited", http.StatusForbidden, `{"message":"API rate limit exceeded for 1.2.3.4."}`, "API rate limit exceeded for 1.2.3.4."},
		{"non-JSON body", http.StatusBadGateway, `<html>bad gateway</html>`, "502 Bad Gateway"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(tt.status)
				_, _ = w.Write([]byte(tt.body))
			}))
			defer srv.Close()

			rel, err := NewClient(srv.Client(), srv.URL).LatestRelease(t.Context(), "o", "r")
			assert.Nil(t, rel)

			var apiErr *APIError
			require.ErrorAs(t, err, &apiErr)
			assert.Equal(t, tt.status, apiErr.StatusCode)
			assert.Equal(t, tt.wantMsg, apiErr.Message)
		})
	}
}
