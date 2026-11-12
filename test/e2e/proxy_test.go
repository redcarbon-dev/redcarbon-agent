//go:build e2e

package e2e

import (
	"io"
	"net/http"
	"net/http/httptest"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/durationpb"

	agents_publicv1 "pkg.redcarbon.ai/proto/redcarbon/agents_public/v1"
)

func TestAgentAnswersProxyRequests(t *testing.T) {
	gh := newFakeGitHub(t)
	api := newFakeRedCarbonAPI(t)
	platform := runtime.GOOS + "_" + runtime.GOARCH

	type vendorCall struct {
		method, auth, body string
	}
	vendorCalls := make(chan vendorCall, 1)
	vendor := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		vendorCalls <- vendorCall{method: r.Method, auth: r.Header.Get("Authorization"), body: string(body)}
		w.Header().Set("X-Vendor", "fake")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"offenses":[]}`))
	}))
	t.Cleanup(vendor.Close)

	binary := buildAgent(t, t.TempDir(), installedVersion, platform, gh.URL)
	gh.publish(installedVersion, nil)
	dir := install(t, binary, api.URL)

	api.enqueueProxyRequest(&agents_publicv1.AgentRequest{
		RequestId: "req-1",
		Method:    http.MethodPost,
		Url:       vendor.URL + "/api/siem/offenses",
		Headers:   map[string]*agents_publicv1.ValueHeader{"Authorization": {Values: []string{"Bearer vendor-token"}}},
		Body:      []byte(`{"filter":"status=OPEN"}`),
		Timeout:   durationpb.New(waitTimeout),
	})

	a := startAgent(t, dir)
	a.waitFor(t, "the proxy response", func() bool { return len(api.receivedProxyResponses()) > 0 })
	a.interrupt(t)

	require.Len(t, vendorCalls, 1)
	assert.Equal(t, vendorCall{method: http.MethodPost, auth: "Bearer vendor-token", body: `{"filter":"status=OPEN"}`}, <-vendorCalls)

	res := api.receivedProxyResponses()[0]
	assert.Equal(t, "req-1", res.GetRequestId())
	assert.Equal(t, int32(http.StatusCreated), res.GetResponse().GetStatus())
	assert.Equal(t, `{"offenses":[]}`, string(res.GetResponse().GetBody()))
	assert.Equal(t, []string{"fake"}, res.GetResponse().GetHeaders()["X-Vendor"].GetValues())
}
