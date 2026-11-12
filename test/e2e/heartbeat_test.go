//go:build e2e

package e2e

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHeartbeatReportsBuildVersion(t *testing.T) {
	api, dir := installLatest(t)

	a := startAgent(t, dir)
	a.waitFor(t, "a heartbeat", func() bool { return len(api.receivedHeartbeats()) > 0 })
	a.interrupt(t)

	hb := api.receivedHeartbeats()[0]
	assert.Equal(t, installedVersion, hb.msg.GetVersion())
	assert.Equal(t, "redcarbon-agent/"+installedVersion, hb.userAgent)
}
