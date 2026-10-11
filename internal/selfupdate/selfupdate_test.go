package selfupdate

import (
	"bytes"
	"errors"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"testing"
	"testing/iotest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func writeTarget(t *testing.T, content string, mode os.FileMode) string {
	t.Helper()

	target := filepath.Join(t.TempDir(), "redcarbon")
	require.NoError(t, os.WriteFile(target, []byte(content), mode))
	require.NoError(t, os.Chmod(target, mode))

	return target
}

func assertOnlyTarget(t *testing.T, target string) {
	t.Helper()

	entries, err := os.ReadDir(filepath.Dir(target))
	require.NoError(t, err)

	var names []string
	for _, e := range entries {
		names = append(names, e.Name())
	}

	assert.Equal(t, []string{filepath.Base(target)}, names)
}

func assertContent(t *testing.T, path, expected string) {
	t.Helper()

	got, err := os.ReadFile(path)
	require.NoError(t, err)
	assert.Equal(t, expected, string(got))
}

func TestReplaceKeepsMode(t *testing.T) {
	target := writeTarget(t, "old", 0o750)

	require.NoError(t, Replace(target, strings.NewReader("new")))

	assertContent(t, target, "new")
	info, err := os.Stat(target)
	require.NoError(t, err)
	assert.Equal(t, os.FileMode(0o750), info.Mode().Perm())
	assertOnlyTarget(t, target)
}

func TestReplaceThroughSymlinkReplacesLinkTarget(t *testing.T) {
	target := writeTarget(t, "old", 0o755)
	link := filepath.Join(t.TempDir(), "redcarbon")
	require.NoError(t, os.Symlink(target, link))

	require.NoError(t, Replace(link, strings.NewReader("new")))

	dest, err := os.Readlink(link)
	require.NoError(t, err)
	assert.Equal(t, target, dest)
	assertContent(t, target, "new")
	assertOnlyTarget(t, target)
}

func TestReplaceKeepsOpenHandlesOnOldContent(t *testing.T) {
	target := writeTarget(t, "old", 0o755)

	running, err := os.Open(target)
	require.NoError(t, err)
	defer running.Close()

	require.NoError(t, Replace(target, strings.NewReader("new")))

	got, err := io.ReadAll(running)
	require.NoError(t, err)
	assert.Equal(t, "old", string(got))
	assertContent(t, target, "new")
}

func TestReplaceConcurrent(t *testing.T) {
	target := writeTarget(t, "old", 0o755)
	payload := bytes.Repeat([]byte("x"), 1<<20)

	var wg sync.WaitGroup
	for range 8 {
		wg.Go(func() {
			assert.NoError(t, Replace(target, bytes.NewReader(payload)))
		})
	}
	wg.Wait()

	got, err := os.ReadFile(target)
	require.NoError(t, err)
	assert.Equal(t, payload, got)
	assertOnlyTarget(t, target)
}

func buildFakeAgent(t *testing.T, version string) string {
	t.Helper()

	goBin, err := exec.LookPath("go")
	if err != nil {
		t.Skip("go toolchain not available")
	}

	out := filepath.Join(t.TempDir(), "fakeagent")
	build := exec.Command(goBin, "build", "-ldflags=-X main.version="+version, "-o", out, "./testdata/fakeagent")
	build.Env = append(os.Environ(), "CGO_ENABLED=0")
	output, err := build.CombinedOutput()
	require.NoError(t, err, string(output))

	return out
}

// On darwin, rewriting a signed Mach-O in place gets the next exec SIGKILLed by the
// kernel's code signature cache; swapping in a new inode avoids it.
func TestReplaceRunningBinary(t *testing.T) {
	if testing.Short() {
		t.Skip("builds binaries")
	}

	v1, err := os.ReadFile(buildFakeAgent(t, "v1"))
	require.NoError(t, err)
	target := writeTarget(t, string(v1), 0o755)

	running := exec.Command(target, "wait")
	stdout, err := running.StdoutPipe()
	require.NoError(t, err)
	require.NoError(t, running.Start())
	t.Cleanup(func() {
		_ = running.Process.Kill()
		_ = running.Wait()
	})

	line := make([]byte, 3)
	_, err = io.ReadFull(stdout, line)
	require.NoError(t, err)
	require.Equal(t, "v1\n", string(line))

	v2, err := os.Open(buildFakeAgent(t, "v2"))
	require.NoError(t, err)
	defer v2.Close()

	require.NoError(t, Replace(target, v2))

	out, err := exec.Command(target).Output()
	require.NoError(t, err)
	assert.Equal(t, "v2\n", string(out))
	assert.NoError(t, running.Process.Signal(syscall.Signal(0)), "the old binary should still be running")
	assertOnlyTarget(t, target)
}

func TestExecutableResolvesSymlinks(t *testing.T) {
	exe, err := Executable()
	require.NoError(t, err)

	resolved, err := filepath.EvalSymlinks(exe)
	require.NoError(t, err)
	assert.Equal(t, resolved, exe)
	assert.True(t, filepath.IsAbs(exe))
}

func TestReplaceFailureLeavesTargetUntouched(t *testing.T) {
	tests := map[string]io.Reader{
		"empty source":            strings.NewReader(""),
		"source fails midway":     io.MultiReader(strings.NewReader("partial"), iotest.ErrReader(errors.New("boom"))),
		"source fails right away": iotest.ErrReader(errors.New("boom")),
	}

	for name, src := range tests {
		t.Run(name, func(t *testing.T) {
			target := writeTarget(t, "old", 0o755)

			require.Error(t, Replace(target, src))

			assertContent(t, target, "old")
			assertOnlyTarget(t, target)
		})
	}
}

func TestReplaceRenameFailureLeavesTargetUntouched(t *testing.T) {
	rename = func(string, string) error { return errors.New("rename failed") }
	t.Cleanup(func() { rename = os.Rename })

	target := writeTarget(t, "old", 0o755)

	require.EqualError(t, Replace(target, strings.NewReader("new")), "rename failed")

	assertContent(t, target, "old")
	assertOnlyTarget(t, target)
}

func TestReplaceUnwritableDirectory(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("root ignores directory permissions")
	}

	target := writeTarget(t, "old", 0o755)
	dir := filepath.Dir(target)
	require.NoError(t, os.Chmod(dir, 0o555))
	t.Cleanup(func() { _ = os.Chmod(dir, 0o755) })

	require.ErrorIs(t, Replace(target, strings.NewReader("new")), os.ErrPermission)

	assertContent(t, target, "old")
	assertOnlyTarget(t, target)
}

func TestReplaceMissingTarget(t *testing.T) {
	dir := t.TempDir()

	require.ErrorIs(t, Replace(filepath.Join(dir, "redcarbon"), strings.NewReader("new")), os.ErrNotExist)

	entries, err := os.ReadDir(dir)
	require.NoError(t, err)
	assert.Empty(t, entries)
}

func TestReplaceRejectsDirectory(t *testing.T) {
	dir := t.TempDir()

	require.Error(t, Replace(dir, strings.NewReader("new")))

	entries, err := os.ReadDir(dir)
	require.NoError(t, err)
	assert.Empty(t, entries)
}
