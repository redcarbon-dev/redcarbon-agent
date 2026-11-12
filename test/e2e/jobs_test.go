//go:build e2e

package e2e

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAgentRunsNoJobs(t *testing.T) {
	gh := newFakeGitHub(t)
	api := newFakeRedCarbonAPI(t)
	platform := runtime.GOOS + "_" + runtime.GOARCH

	binary := buildAgent(t, t.TempDir(), installedVersion, platform, gh.URL)
	gh.publish(installedVersion, nil)
	dir := install(t, binary, api.URL)

	a := startAgent(t, dir)
	a.waitFor(t, "two heartbeats", func() bool { return len(api.receivedHeartbeats()) >= 2 })
	a.interrupt(t)

	assert.NotContains(t, api.calledProcedures(), "/redcarbon.agents_public.v1.AgentsPublicAPIsV1Srv/FetchAgentConfiguration")
	assert.NotContains(t, api.calledProcedures(), "/redcarbon.agents_public.v1.AgentsPublicAPIsV1Srv/IngestIncident")
}
