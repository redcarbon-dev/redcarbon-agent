//go:build e2e

package e2e

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/fs"
	"maps"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"sync"
	"syscall"
	"testing"
	"time"

	"connectrpc.com/connect"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"

	agents_publicv1 "pkg.redcarbon.ai/proto/redcarbon/agents_public/v1"
	"pkg.redcarbon.ai/proto/redcarbon/agents_public/v1/agents_publicv1connect"
)

const (
	installedVersion = "0.0.1-e2e.1"
	releasedVersion  = "0.0.2-e2e.1"

	updateExitCode = 3
	installedMode  = fs.FileMode(0o750)

	waitTimeout = time.Minute
	gracePeriod = 2 * time.Second
)

func TestSelfUpdate(t *testing.T) {
	gh := newFakeGitHub(t)
	apiURL := newFakeRedCarbonAPI(t)
	platform := runtime.GOOS + "_" + runtime.GOARCH

	buildDir := t.TempDir()
	installed := buildAgent(t, buildDir, installedVersion, platform, gh.URL)
	released := buildAgent(t, buildDir, releasedVersion, platform, gh.URL)
	archive := buildArchive(t, releasedVersion, platform, released)
	archivePath := "/download/" + archiveName(releasedVersion, platform)

	goodRelease := func() map[string][]byte {
		assets := map[string][]byte{}
		sums := map[string][]byte{}
		for _, p := range releasePlatforms {
			name := archiveName(releasedVersion, p)
			assets[name] = nil
			sums[name] = sha256Sum([]byte(p))
		}
		assets[archiveName(releasedVersion, platform)] = archive
		sums[archiveName(releasedVersion, platform)] = sha256Sum(archive)
		assets[checksumsName(releasedVersion)] = checksumsFile(sums)
		return assets
	}

	t.Run("success", func(t *testing.T) {
		gh.publish(releasedVersion, goodRelease())
		dir := install(t, installed, apiURL)

		a := startAgent(t, dir)
		require.Equal(t, updateExitCode, a.waitExit(t), a.output())

		version, err := exec.Command(filepath.Join(dir, "bin", "redcarbon"), "--version").CombinedOutput()
		require.NoError(t, err, string(version))
		assert.Contains(t, string(version), releasedVersion)

		info, err := os.Stat(filepath.Join(dir, "bin", "redcarbon"))
		require.NoError(t, err)
		assert.Equal(t, installedMode, info.Mode().Perm())
		assert.Equal(t, []string{"redcarbon"}, listDir(t, filepath.Join(dir, "bin")))
		require.Equal(t, 1, gh.hits(archivePath))

		restarted := startAgent(t, dir)
		restarted.waitFor(t, "latest release check after restart", func() bool {
			return gh.hits(latestReleasePath) == 2 && restarted.logged("already at the last version")
		})
		restarted.requireStaysAlive(t)
		assert.Equal(t, 1, gh.hits(archivePath), "the updated agent must not download the release again")
		restarted.interrupt(t)
	})

	t.Run("wrong checksum", func(t *testing.T) {
		assets := goodRelease()
		assets[checksumsName(releasedVersion)] = checksumsFile(map[string][]byte{
			archiveName(releasedVersion, platform): sha256Sum([]byte("tampered")),
		})
		gh.publish(releasedVersion, assets)
		dir := install(t, installed, apiURL)

		a := startAgent(t, dir)
		a.waitFor(t, "checksum mismatch", func() bool {
			return gh.hits(archivePath) == 1 && a.logged("wrong checksum")
		})
		a.requireStaysAlive(t)
		requireInstalled(t, dir, installed)
		a.interrupt(t)
	})

	t.Run("read-only bin directory", func(t *testing.T) {
		if os.Geteuid() == 0 {
			t.Skip("root ignores directory permissions")
		}

		gh.publish(releasedVersion, goodRelease())
		dir := install(t, installed, apiURL)
		binDir := filepath.Join(dir, "bin")
		require.NoError(t, os.Chmod(binDir, 0o555))
		t.Cleanup(func() { _ = os.Chmod(binDir, 0o755) })

		a := startAgent(t, dir)
		a.waitFor(t, "failed replacement", func() bool {
			return gh.hits(archivePath) == 1 && a.logged("Unexpected error while updating the binary") && a.logged("permission denied")
		})
		a.requireStaysAlive(t)
		requireInstalled(t, dir, installed)
		a.interrupt(t)
	})

	t.Run("archive not found", func(t *testing.T) {
		assets := goodRelease()
		assets[archiveName(releasedVersion, platform)] = nil
		gh.publish(releasedVersion, assets)
		dir := install(t, installed, apiURL)

		a := startAgent(t, dir)
		a.waitFor(t, "failed download", func() bool {
			return gh.hits(archivePath) == 1 && a.logged("unexpected status 404")
		})
		a.requireStaysAlive(t)
		requireInstalled(t, dir, installed)
		a.interrupt(t)
	})
}

func buildAgent(t *testing.T, outDir, version, platform, githubURL string) []byte {
	t.Helper()

	out := filepath.Join(outDir, "redcarbon-"+version)
	ldflags := strings.Join([]string{
		"-X pkg.redcarbon.ai/internal/build.Version=" + version,
		"-X pkg.redcarbon.ai/internal/build.Architecture=" + platform,
		"-X pkg.redcarbon.ai/internal/build.GitHubAPIURL=" + githubURL,
	}, " ")

	cmd := exec.Command("go", "build", "-ldflags", ldflags, "-o", out, "./cmd")
	cmd.Dir = repoRoot
	cmd.Env = append(os.Environ(), "CGO_ENABLED=0")

	b, err := cmd.CombinedOutput()
	require.NoError(t, err, string(b))

	// -X silently ignores unknown symbols: without this check a rename would point the agent at the real GitHub.
	symbols, err := exec.Command("go", "tool", "nm", out).CombinedOutput()
	require.NoError(t, err, string(symbols))
	require.Contains(t, string(symbols), " pkg.redcarbon.ai/internal/build.GitHubAPIURL\n")

	binary, err := os.ReadFile(out)
	require.NoError(t, err)

	return binary
}

// install mimics an extracted release whose agent is configured with a single profile.
func install(t *testing.T, binary []byte, apiURL string) string {
	t.Helper()

	dir := t.TempDir()
	require.NoError(t, os.Mkdir(filepath.Join(dir, "bin"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "bin", "redcarbon"), binary, installedMode))
	require.NoError(t, os.Chmod(filepath.Join(dir, "bin", "redcarbon"), installedMode))

	config := fmt.Sprintf("version: \"1.0\"\nprofiles:\n    - name: e2e\n      profile:\n        host: %s\n        token: e2e-token\n", apiURL)
	require.NoError(t, os.WriteFile(filepath.Join(dir, ".config.yaml"), []byte(config), 0o600))

	return dir
}

func requireInstalled(t *testing.T, dir string, binary []byte) {
	t.Helper()

	b, err := os.ReadFile(filepath.Join(dir, "bin", "redcarbon"))
	require.NoError(t, err)
	require.Equal(t, sha256.Sum256(binary), sha256.Sum256(b), "the installed binary must be left untouched")
	require.Equal(t, []string{"redcarbon"}, listDir(t, filepath.Join(dir, "bin")))
}

func listDir(t *testing.T, dir string) []string {
	t.Helper()

	entries, err := os.ReadDir(dir)
	require.NoError(t, err)

	var names []string
	for _, e := range entries {
		names = append(names, e.Name())
	}

	return names
}

// agent is a running agent process; the test plays the supervisor (systemd Restart=on-failure).
type agent struct {
	cmd    *exec.Cmd
	out    *syncBuffer
	exited chan struct{}
}

func startAgent(t *testing.T, dir string) *agent {
	t.Helper()

	out := &syncBuffer{}
	cmd := exec.Command(filepath.Join(dir, "bin", "redcarbon"))
	cmd.Dir = dir
	cmd.Stdout = out
	cmd.Stderr = out
	require.NoError(t, cmd.Start())

	a := &agent{cmd: cmd, out: out, exited: make(chan struct{})}
	go func() {
		_ = cmd.Wait()
		close(a.exited)
	}()

	t.Cleanup(func() {
		if a.running() {
			_ = cmd.Process.Kill()
			<-a.exited
		}
		if t.Failed() {
			t.Logf("agent output:\n%s", a.output())
		}
	})

	return a
}

func (a *agent) running() bool {
	select {
	case <-a.exited:
		return false
	default:
		return true
	}
}

func (a *agent) output() string {
	return a.out.String()
}

func (a *agent) logged(s string) bool {
	return strings.Contains(a.output(), s)
}

func (a *agent) waitExit(t *testing.T) int {
	t.Helper()

	select {
	case <-a.exited:
		return a.cmd.ProcessState.ExitCode()
	case <-time.After(waitTimeout):
		require.FailNow(t, "agent did not exit")
		return 0
	}
}

// waitFor polls cond and fails as soon as the agent exits, so a regression that
// restarts the agent surfaces immediately instead of as a timeout.
func (a *agent) waitFor(t *testing.T, what string, cond func() bool) {
	t.Helper()

	deadline := time.Now().Add(waitTimeout)
	for !cond() {
		if !a.running() {
			require.FailNow(t, "agent exited while waiting for "+what, "exit code %d", a.cmd.ProcessState.ExitCode())
		}
		if time.Now().After(deadline) {
			require.FailNow(t, "timed out waiting for "+what)
		}
		time.Sleep(50 * time.Millisecond)
	}
}

func (a *agent) requireStaysAlive(t *testing.T) {
	t.Helper()

	select {
	case <-a.exited:
		require.FailNow(t, "agent exited", "exit code %d", a.cmd.ProcessState.ExitCode())
	case <-time.After(gracePeriod):
	}
}

func (a *agent) interrupt(t *testing.T) {
	t.Helper()

	require.NoError(t, a.cmd.Process.Signal(syscall.SIGINT))
	require.Equal(t, 0, a.waitExit(t))
}

type syncBuffer struct {
	mu  sync.Mutex
	buf bytes.Buffer
}

func (b *syncBuffer) Write(p []byte) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buf.Write(p)
}

func (b *syncBuffer) String() string {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buf.String()
}

const latestReleasePath = "/repos/redcarbon-dev/redcarbon-agent/releases/latest"

// fakeGitHub serves the latest release and its assets; a nil asset is listed but answers 404.
type fakeGitHub struct {
	*httptest.Server

	mu       sync.Mutex
	tag      string
	assets   map[string][]byte
	requests map[string]int
}

func newFakeGitHub(t *testing.T) *fakeGitHub {
	gh := &fakeGitHub{requests: map[string]int{}}
	gh.Server = httptest.NewServer(http.HandlerFunc(gh.serveHTTP))
	t.Cleanup(gh.Close)

	return gh
}

func (gh *fakeGitHub) publish(version string, assets map[string][]byte) {
	gh.mu.Lock()
	defer gh.mu.Unlock()

	gh.tag = "v" + version
	gh.assets = assets
	gh.requests = map[string]int{}
}

func (gh *fakeGitHub) hits(path string) int {
	gh.mu.Lock()
	defer gh.mu.Unlock()

	return gh.requests[path]
}

func (gh *fakeGitHub) serveHTTP(w http.ResponseWriter, r *http.Request) {
	gh.mu.Lock()
	defer gh.mu.Unlock()

	gh.requests[r.URL.Path]++

	if r.URL.Path == latestReleasePath {
		type asset struct {
			Name               string `json:"name"`
			BrowserDownloadURL string `json:"browser_download_url"`
		}

		var assets []asset
		for _, name := range slices.Sorted(maps.Keys(gh.assets)) {
			assets = append(assets, asset{Name: name, BrowserDownloadURL: gh.URL + "/download/" + name})
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{"tag_name": gh.tag, "assets": assets})

		return
	}

	content, ok := gh.assets[strings.TrimPrefix(r.URL.Path, "/download/")]
	if !ok || content == nil {
		http.NotFound(w, r)
		return
	}

	_, _ = w.Write(content)
}

type fakeRedCarbonAPI struct {
	agents_publicv1connect.UnimplementedAgentsPublicAPIsV1SrvHandler
}

func newFakeRedCarbonAPI(t *testing.T) string {
	mux := http.NewServeMux()
	mux.Handle(agents_publicv1connect.NewAgentsPublicAPIsV1SrvHandler(fakeRedCarbonAPI{}))

	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)

	return srv.URL
}

func (fakeRedCarbonAPI) HZ(context.Context, *connect.Request[agents_publicv1.HZRequest]) (*connect.Response[agents_publicv1.HZResponse], error) {
	return connect.NewResponse(&agents_publicv1.HZResponse{ReceivedAt: timestamppb.Now()}), nil
}

func (fakeRedCarbonAPI) FetchAgentConfiguration(context.Context, *connect.Request[agents_publicv1.FetchAgentConfigurationRequest]) (*connect.Response[agents_publicv1.FetchAgentConfigurationResponse], error) {
	return connect.NewResponse(&agents_publicv1.FetchAgentConfigurationResponse{}), nil
}

func (fakeRedCarbonAPI) FetchAgentRequests(context.Context, *connect.Request[agents_publicv1.FetchAgentRequestsRequest]) (*connect.Response[agents_publicv1.FetchAgentRequestsResponse], error) {
	return connect.NewResponse(&agents_publicv1.FetchAgentRequestsResponse{}), nil
}
