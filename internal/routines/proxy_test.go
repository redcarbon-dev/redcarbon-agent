package routines

import (
	"bytes"
	"context"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"

	"connectrpc.com/connect"
	"pkg.redcarbon.ai/internal/config"
	agents_publicv1 "pkg.redcarbon.ai/proto/redcarbon/agents_public/v1"
	"pkg.redcarbon.ai/proto/redcarbon/agents_public/v1/agents_publicv1connect"
)

type stubAgentsClient struct {
	agents_publicv1connect.AgentsPublicAPIsV1SrvClient
}

func (stubAgentsClient) SubmitAgentResponse(context.Context, *connect.Request[agents_publicv1.SubmitAgentResponseRequest]) (*connect.Response[agents_publicv1.SubmitAgentResponseResponse], error) {
	return connect.NewResponse(&agents_publicv1.SubmitAgentResponseResponse{}), nil
}

func TestProcessRequestLogKeepsCredentialsOut(t *testing.T) {
	var logs bytes.Buffer
	prev := slog.Default()
	slog.SetDefault(slog.New(slog.NewTextHandler(&logs, &slog.HandlerOptions{Level: slog.LevelDebug})))
	t.Cleanup(func() { slog.SetDefault(prev) })

	vendor := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	defer vendor.Close()

	vendorURL, err := url.Parse(vendor.URL)
	if err != nil {
		t.Fatal(err)
	}
	vendorURL.User = url.UserPassword("user", "userinfo-secret")
	vendorURL.Path = "/api"
	vendorURL.RawQuery = "apikey=query-secret"

	r := RoutineConfig{
		agentsCli: stubAgentsClient{},
		profile:   config.Profile{Name: "acme"},
	}
	r.processRequest(t.Context(), &agents_publicv1.AgentRequest{
		RequestId: "req-1",
		Method:    http.MethodGet,
		Url:       vendorURL.String(),
		Headers: map[string]*agents_publicv1.ValueHeader{
			"Authorization": {Values: []string{"Basic header-secret"}},
		},
	})

	out := logs.String()
	for _, secret := range []string{"userinfo-secret", "query-secret", "header-secret"} {
		if strings.Contains(out, secret) {
			t.Errorf("log leaks %q:\n%s", secret, out)
		}
	}
	for _, want := range []string{
		`msg="Executing HTTP request"`,
		"profile=acme",
		"apikey=xxxxx",
		"headers=[Authorization]",
		"status=200",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("log missing %q:\n%s", want, out)
		}
	}
	for line := range strings.Lines(out) {
		if strings.Contains(line, `msg="Executing HTTP request"`) && !strings.Contains(line, "id=req-1") {
			t.Errorf("request log line lacks the request id: %s", line)
		}
	}
}

func TestHeaderNamesOmitsValues(t *testing.T) {
	h := http.Header{}
	h.Set("Authorization", "Basic c2VjcmV0")
	h.Set("X-Api-Key", "secret")
	h.Set("Content-Type", "application/json")

	got := headerNames(h)

	want := []string{"Authorization", "Content-Type", "X-Api-Key"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("headerNames() = %v, want %v", got, want)
	}
}

func TestRedactedURLMasksSecrets(t *testing.T) {
	tests := []struct {
		name string
		raw  string
		want string
	}{
		{"no query", "https://wazuh:9200/wazuh-alerts-*/_search", "https://wazuh:9200/wazuh-alerts-*/_search"},
		{"query values", "https://vendor.example/api?apikey=secret&page=2", "https://vendor.example/api?apikey=xxxxx&page=xxxxx"},
		{"userinfo password", "https://user:secret@vendor.example/api", "https://user:xxxxx@vendor.example/api"},
		{"unparseable", "http://[::1", "<unparseable URL>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := redactedURL(tt.raw); got != tt.want {
				t.Fatalf("redactedURL(%q) = %q, want %q", tt.raw, got, tt.want)
			}
		})
	}
}
