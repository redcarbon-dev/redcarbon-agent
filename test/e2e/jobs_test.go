//go:build e2e

package e2e

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAgentRunsNoJobs(t *testing.T) {
	api, dir := installLatest(t)

	a := startAgent(t, dir)
	a.waitFor(t, "two heartbeats", func() bool { return len(api.receivedHeartbeats()) >= 2 })
	a.interrupt(t)

	assert.NotContains(t, api.calledProcedures(), "/redcarbon.agents_public.v1.AgentsPublicAPIsV1Srv/FetchAgentConfiguration")
	assert.NotContains(t, api.calledProcedures(), "/redcarbon.agents_public.v1.AgentsPublicAPIsV1Srv/IngestIncident")
}
