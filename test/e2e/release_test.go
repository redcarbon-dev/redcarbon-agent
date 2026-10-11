// Package e2e exercises the agent binary end to end. The self-update scenarios need
// the e2e build tag (see test/e2e/Dockerfile); the goreleaser drift guard always runs.
package e2e

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"maps"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.yaml.in/yaml/v3"
)

const repoRoot = "../.."

// releasePlatforms are the goos_goarch pairs .goreleaser.yml publishes.
var releasePlatforms = []string{"darwin_amd64", "darwin_arm64", "linux_amd64", "linux_arm64"}

func archiveDir(version, platform string) string {
	return fmt.Sprintf("redcarbon_%s_%s", version, platform)
}

func archiveName(version, platform string) string {
	return archiveDir(version, platform) + ".tar.gz"
}

func checksumsName(version string) string {
	return fmt.Sprintf("redcarbon_%s_checksums.txt", version)
}

// buildArchive lays out the archive the way goreleaser does: everything wrapped in a
// directory named after the archive, assets/* at its root, templates/* under bin/templates.
func buildArchive(t *testing.T, version, platform string, binary []byte) []byte {
	t.Helper()

	dir := archiveDir(version, platform)
	entries := map[string][]byte{path.Join(dir, "bin", "redcarbon"): binary}

	for src, dst := range map[string]string{"assets/*": ".", "templates/*": "bin/templates"} {
		files, err := filepath.Glob(filepath.Join(repoRoot, src))
		require.NoError(t, err)
		require.NotEmpty(t, files, src)

		for _, f := range files {
			b, err := os.ReadFile(f)
			require.NoError(t, err)
			entries[path.Join(dir, dst, filepath.Base(f))] = b
		}
	}

	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gz)

	for _, name := range slices.Sorted(maps.Keys(entries)) {
		require.NoError(t, tw.WriteHeader(&tar.Header{Name: name, Mode: 0o755, Size: int64(len(entries[name])), Typeflag: tar.TypeReg}))
		_, err := tw.Write(entries[name])
		require.NoError(t, err)
	}

	require.NoError(t, tw.Close())
	require.NoError(t, gz.Close())

	return buf.Bytes()
}

// checksumsFile renders goreleaser's default sha256 checksums file.
func checksumsFile(sums map[string][]byte) []byte {
	var b strings.Builder

	for _, name := range slices.Sorted(maps.Keys(sums)) {
		fmt.Fprintf(&b, "%s  %s\n", hex.EncodeToString(sums[name]), name)
	}

	return []byte(b.String())
}

func sha256Sum(b []byte) []byte {
	s := sha256.Sum256(b)
	return s[:]
}

type goreleaserConfig struct {
	ProjectName string `yaml:"project_name"`
	Builds      []struct {
		Binary  string   `yaml:"binary"`
		Main    string   `yaml:"main"`
		Ldflags []string `yaml:"ldflags"`
		Goos    []string `yaml:"goos"`
		Goarch  []string `yaml:"goarch"`
	} `yaml:"builds"`
	Archives []struct {
		NameTemplate    string   `yaml:"name_template"`
		WrapInDirectory bool     `yaml:"wrap_in_directory"`
		Formats         []string `yaml:"formats"`
		Files           []struct {
			Src string `yaml:"src"`
			Dst string `yaml:"dst"`
		} `yaml:"files"`
	} `yaml:"archives"`
	Checksum map[string]any `yaml:"checksum"`
}

// TestGoreleaserConfigMatchesEmulatedRelease fails when .goreleaser.yml drifts away from
// the release layout the self-update e2e test emulates (and the update routine expects).
func TestGoreleaserConfigMatchesEmulatedRelease(t *testing.T) {
	b, err := os.ReadFile(filepath.Join(repoRoot, ".goreleaser.yml"))
	require.NoError(t, err)

	var c goreleaserConfig
	require.NoError(t, yaml.Unmarshal(b, &c))

	assert.Equal(t, "redcarbon", c.ProjectName)

	var platforms []string
	for _, build := range c.Builds {
		assert.Equal(t, "bin/redcarbon", build.Binary)
		assert.Equal(t, "./cmd", build.Main)

		ldflags := strings.Join(build.Ldflags, " ")
		assert.Contains(t, ldflags, "-X pkg.redcarbon.ai/internal/build.Version={{.Version}} ")
		assert.Contains(t, ldflags, "-X pkg.redcarbon.ai/internal/build.Architecture={{ .Os }}_{{ .Arch }}{{ if .Arm }}v{{ .Arm }}{{ end }} ")
		assert.NotContains(t, ldflags, "GitHubAPIURL", "releases must talk to the real GitHub API")

		for _, goos := range build.Goos {
			for _, goarch := range build.Goarch {
				platforms = append(platforms, goos+"_"+goarch)
			}
		}
	}
	slices.Sort(platforms)
	assert.Equal(t, releasePlatforms, platforms)

	require.Len(t, c.Archives, 1)
	archive := c.Archives[0]
	assert.Equal(t, "{{ .ProjectName }}_{{ .Version }}_{{ .Os }}_{{ .Arch }}{{ if .Arm }}v{{ .Arm }}{{ end }}", archive.NameTemplate)
	assert.True(t, archive.WrapInDirectory)
	assert.Equal(t, []string{"tar.gz"}, archive.Formats)

	files := map[string]string{}
	for _, f := range archive.Files {
		files[f.Src] = f.Dst
	}
	assert.Equal(t, map[string]string{"assets/*": ".", "templates/*": "bin/templates"}, files)

	assert.Nil(t, c.Checksum, "the e2e test emulates goreleaser's default checksum file name and sha256 algorithm")
}
