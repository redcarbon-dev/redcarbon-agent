package routines

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"maps"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"slices"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"pkg.redcarbon.ai/internal/build"
	"pkg.redcarbon.ai/internal/config"
	"pkg.redcarbon.ai/internal/github"
)

const (
	archiveName   = "redcarbon_1.2.3_linux_amd64.tar.gz"
	checksumsName = "redcarbon_1.2.3_checksums.txt"
	newBinary     = "new-binary"
)

type releaseAsset struct {
	name    string
	handler http.HandlerFunc
}

type releaseServer struct {
	srv      *httptest.Server
	requests atomic.Int32
}

func newReleaseServer(t *testing.T, assets ...releaseAsset) *releaseServer {
	t.Helper()

	rs := &releaseServer{}
	mux := http.NewServeMux()

	mux.HandleFunc("GET /repos/redcarbon-dev/redcarbon-agent/releases/latest", func(w http.ResponseWriter, _ *http.Request) {
		var list []map[string]string
		for _, a := range assets {
			list = append(list, map[string]string{"name": a.name, "browser_download_url": rs.srv.URL + "/download/" + a.name})
		}

		_ = json.NewEncoder(w).Encode(map[string]any{"tag_name": "v1.2.3", "assets": list})
	})

	for _, a := range assets {
		mux.HandleFunc("GET /download/"+a.name, a.handler)
	}

	rs.srv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rs.requests.Add(1)
		assert.Equal(t, "redcarbon-agent/"+build.Version, r.Header.Get("User-Agent"), r.URL.Path)
		mux.ServeHTTP(w, r)
	}))
	t.Cleanup(rs.srv.Close)

	return rs
}

func (rs *releaseServer) routine(done chan bool) RoutineConfig {
	return RoutineConfig{
		gh:      github.NewClient(rs.srv.Client(), rs.srv.URL),
		done:    done,
		profile: config.Profile{Name: "acme"},
	}
}

func serve(content []byte) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write(content)
	}
}

func buildArchive(t *testing.T, files map[string]string) []byte {
	t.Helper()

	var buf bytes.Buffer
	gzw := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gzw)

	require.NoError(t, tw.WriteHeader(&tar.Header{Name: "redcarbon_1.2.3_linux_amd64/", Typeflag: tar.TypeDir, Mode: 0o755}))

	for _, name := range slices.Sorted(maps.Keys(files)) {
		content := files[name]
		require.NoError(t, tw.WriteHeader(&tar.Header{Name: name, Typeflag: tar.TypeReg, Mode: 0o755, Size: int64(len(content))}))
		_, err := tw.Write([]byte(content))
		require.NoError(t, err)
	}

	require.NoError(t, tw.Close())
	require.NoError(t, gzw.Close())

	return buf.Bytes()
}

func releaseArchive(t *testing.T) []byte {
	t.Helper()

	return buildArchive(t, map[string]string{
		"redcarbon_1.2.3_linux_amd64/svc.sh":        "#!/bin/bash",
		"redcarbon_1.2.3_linux_amd64/bin/redcarbon": newBinary,
	})
}

func sha256Hex(content []byte) string {
	sum := sha256.Sum256(content)
	return hex.EncodeToString(sum[:])
}

func validRelease(t *testing.T) []releaseAsset {
	t.Helper()

	archive := releaseArchive(t)

	return []releaseAsset{
		{name: archiveName, handler: serve(archive)},
		{name: checksumsName, handler: serve([]byte(sha256Hex(archive) + "  " + archiveName + "\n"))},
	}
}

// installAgent pretends the running agent is version 1.0.0 for linux_amd64 installed at
// the returned path, which holds "old-binary".
func installAgent(t *testing.T) string {
	t.Helper()

	target := filepath.Join(t.TempDir(), "redcarbon")
	require.NoError(t, os.WriteFile(target, []byte("old-binary"), 0o755))

	prevVersion, prevArch, prevExecutable := build.Version, build.Architecture, executable
	build.Version, build.Architecture = "1.0.0", "linux_amd64"
	executable = func() (string, error) { return target, nil }
	t.Cleanup(func() {
		build.Version, build.Architecture, executable = prevVersion, prevArch, prevExecutable
	})

	return target
}

func assertBinary(t *testing.T, target, expected string) {
	t.Helper()

	got, err := os.ReadFile(target)
	require.NoError(t, err)
	assert.Equal(t, expected, string(got))
}

func TestUpdateRoutineReplacesExecutableAndSignalsDone(t *testing.T) {
	target := installAgent(t)
	done := make(chan bool, 1)

	newReleaseServer(t, validRelease(t)...).routine(done).UpdateRoutine(t.Context())

	assertBinary(t, target, newBinary)
	assert.Len(t, done, 1)
}

func TestUpdateRoutineDoesNotSignalDoneWhenReplaceFails(t *testing.T) {
	tests := map[string]func(t *testing.T, target string){
		"target is not a regular file": func(t *testing.T, target string) {
			require.NoError(t, os.Remove(target))
			require.NoError(t, os.Mkdir(target, 0o755))
		},
		"install directory is not writable": func(t *testing.T, target string) {
			if os.Geteuid() == 0 {
				t.Skip("root ignores directory permissions")
			}

			dir := filepath.Dir(target)
			require.NoError(t, os.Chmod(dir, 0o555))
			t.Cleanup(func() { _ = os.Chmod(dir, 0o755) })
		},
	}

	for name, breakTarget := range tests {
		t.Run(name, func(t *testing.T) {
			target := installAgent(t)
			breakTarget(t, target)
			done := make(chan bool, 1)

			newReleaseServer(t, validRelease(t)...).routine(done).UpdateRoutine(t.Context())

			assert.Empty(t, done, "a failed replace must not restart the agent")
		})
	}
}

func TestUpdateRoutineOnlyOneProfileUpdatesAtATime(t *testing.T) {
	target := installAgent(t)
	archive := releaseArchive(t)
	downloading := make(chan struct{})
	startedDownloading := sync.OnceFunc(func() { close(downloading) })
	release := make(chan struct{})

	rs := newReleaseServer(t,
		releaseAsset{name: archiveName, handler: func(w http.ResponseWriter, _ *http.Request) {
			startedDownloading()
			<-release
			_, _ = w.Write(archive)
		}},
		releaseAsset{name: checksumsName, handler: serve([]byte(sha256Hex(archive) + "  " + archiveName + "\n"))},
	)

	// Unbuffered like in main: a second profile reaching the send would block forever.
	done := make(chan bool)
	first := make(chan struct{})

	go func() {
		defer close(first)
		rs.routine(done).UpdateRoutine(t.Context())
	}()

	<-downloading
	requestsBefore := rs.requests.Load()

	second := rs.routine(done)
	second.profile.Name = "other"
	second.UpdateRoutine(t.Context())

	assert.Equal(t, requestsBefore, rs.requests.Load(), "the second profile must not touch the release while the first is updating")

	close(release)
	<-done
	<-first

	assertBinary(t, target, newBinary)
}

func TestUpdateRoutineDoesNotBlockOnDoneAfterShutdown(t *testing.T) {
	target := installAgent(t)
	ctx, cancel := context.WithCancel(t.Context())
	rs := newReleaseServer(t, validRelease(t)...)

	finished := make(chan struct{})
	go func() {
		defer close(finished)
		rs.routine(make(chan bool)).UpdateRoutine(ctx)
	}()

	require.Eventually(t, func() bool {
		got, _ := os.ReadFile(target)
		return string(got) == newBinary
	}, 5*time.Second, time.Millisecond)
	cancel()

	select {
	case <-finished:
	case <-time.After(5 * time.Second):
		t.Fatal("UpdateRoutine stayed blocked on done after its context was cancelled")
	}
}

func TestUpdateRoutineMatchesChecksumByExactAssetName(t *testing.T) {
	archive := releaseArchive(t)
	sbom := []byte(`{"sbom": true}`)
	wrong := sha256Hex([]byte("something else"))

	tests := map[string]struct {
		checksums string
		updated   bool
	}{
		"exact line after a longer name sharing the prefix": {
			checksums: wrong + "  " + archiveName + ".sbom.json\n" + sha256Hex(archive) + "  " + archiveName + "\n",
			updated:   true,
		},
		"only a longer name sharing the prefix": {
			checksums: sha256Hex(sbom) + "  " + archiveName + ".sbom.json\n",
		},
		"no line for the asset": {
			checksums: sha256Hex(archive) + "  redcarbon_1.2.3_linux_arm64.tar.gz\n",
		},
		"tampered archive": {
			checksums: wrong + "  " + archiveName + "\n",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			target := installAgent(t)
			done := make(chan bool, 1)

			newReleaseServer(t,
				releaseAsset{name: archiveName + ".sbom.json", handler: serve(sbom)},
				releaseAsset{name: archiveName, handler: serve(archive)},
				releaseAsset{name: checksumsName, handler: serve([]byte(tt.checksums))},
			).routine(done).UpdateRoutine(t.Context())

			if tt.updated {
				assertBinary(t, target, newBinary)
				assert.Len(t, done, 1)
			} else {
				assertBinary(t, target, "old-binary")
				assert.Empty(t, done)
			}
		})
	}
}

func TestUpdateRoutineRejectsFailedDownloads(t *testing.T) {
	archive := releaseArchive(t)
	checksums := []byte(sha256Hex(archive) + "  " + archiveName + "\n")
	status := func(code int, body []byte) http.HandlerFunc {
		return func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(code)
			_, _ = w.Write(body)
		}
	}

	tests := map[string][]releaseAsset{
		"archive not found": {
			{name: archiveName, handler: status(http.StatusNotFound, archive)},
			{name: checksumsName, handler: serve(checksums)},
		},
		"archive server error": {
			{name: archiveName, handler: status(http.StatusInternalServerError, archive)},
			{name: checksumsName, handler: serve(checksums)},
		},
		"checksums server error": {
			{name: archiveName, handler: serve(archive)},
			{name: checksumsName, handler: status(http.StatusBadGateway, checksums)},
		},
	}

	for name, assets := range tests {
		t.Run(name, func(t *testing.T) {
			target := installAgent(t)
			done := make(chan bool, 1)

			newReleaseServer(t, assets...).routine(done).UpdateRoutine(t.Context())

			assertBinary(t, target, "old-binary")
			assert.Empty(t, done)
		})
	}
}

func TestUpdateRoutineCapsDownloadSize(t *testing.T) {
	archive := releaseArchive(t)

	tests := map[string]struct {
		limit   int64
		updated bool
	}{
		"archive at the limit":       {limit: int64(len(archive)), updated: true},
		"archive one byte too large": {limit: int64(len(archive)) - 1},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			target := installAgent(t)
			prev := maxDownloadSize
			maxDownloadSize = tt.limit
			t.Cleanup(func() { maxDownloadSize = prev })
			done := make(chan bool, 1)

			newReleaseServer(t, validRelease(t)...).routine(done).UpdateRoutine(t.Context())

			if tt.updated {
				assertBinary(t, target, newBinary)
				assert.Len(t, done, 1)
			} else {
				assertBinary(t, target, "old-binary")
				assert.Empty(t, done)
			}
		})
	}
}

func TestUpdateRoutineRejectsArchiveWithoutBinary(t *testing.T) {
	target := installAgent(t)
	archive := buildArchive(t, map[string]string{"redcarbon_1.2.3_linux_amd64/svc.sh": "#!/bin/bash"})
	done := make(chan bool, 1)

	newReleaseServer(t,
		releaseAsset{name: archiveName, handler: serve(archive)},
		releaseAsset{name: checksumsName, handler: serve([]byte(sha256Hex(archive) + "  " + archiveName + "\n"))},
	).routine(done).UpdateRoutine(t.Context())

	assertBinary(t, target, "old-binary")
	assert.Empty(t, done)
}

func TestUpdateRoutineSkipsWhenNotNewer(t *testing.T) {
	for _, version := range []string{"1.2.3", "2.0.0", "DEV", "1.3.0-SNAPSHOT-abc"} {
		t.Run(version, func(t *testing.T) {
			target := installAgent(t)
			build.Version = version
			done := make(chan bool, 1)
			rs := newReleaseServer(t, validRelease(t)...)

			rs.routine(done).UpdateRoutine(t.Context())

			assertBinary(t, target, "old-binary")
			assert.Empty(t, done)
			assert.Equal(t, int32(1), rs.requests.Load(), "only the release metadata may be fetched")
		})
	}
}
