//go:build e2e

package e2e

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHeartbeatReportsBuildVersion(t *testing.T) {
	gh := newFakeGitHub(t)
	api := newFakeRedCarbonAPI(t)
	platform := runtime.GOOS + "_" + runtime.GOARCH

	binary := buildAgent(t, t.TempDir(), installedVersion, platform, gh.URL)
	gh.publish(installedVersion, nil)
	dir := install(t, binary, api.URL)

	a := startAgent(t, dir)
	a.waitFor(t, "a heartbeat", func() bool { return len(api.receivedHeartbeats()) > 0 })
	a.interrupt(t)

	hb := api.receivedHeartbeats()[0]
	assert.Equal(t, installedVersion, hb.msg.GetVersion())
	assert.Equal(t, "redcarbon-agent/"+installedVersion, hb.userAgent)
}
