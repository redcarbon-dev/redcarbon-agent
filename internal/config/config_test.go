package config

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.yaml.in/yaml/v3"
)

func readString(t *testing.T, name string) string {
	t.Helper()
	b, err := os.ReadFile(name)
	require.NoError(t, err)
	return string(b)
}

func setupFile(t *testing.T, content []byte) string {
	t.Helper()
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, fileName), content, 0o644))
	require.NoError(t, Init(dir))
	return filepath.Join(dir, fileName)
}

type viperResult struct {
	File      string          `json:"file"`
	Config    json.RawMessage `json:"config"`
	Panic     string          `json:"panic"`
	FileAfter string          `json:"file_after"`
}

// Expected outputs record what the viper-based implementation loaded from, and
// wrote back to, each golden input, restricted to the fields of the current
// profile model. Panic is set where viper crashed or exited.
// LoadConfiguration exits the process on error, so load is called directly.
func TestLoadMatchesViper(t *testing.T) {
	deviations := map[string]func(t *testing.T, c Configuration, err error, input, after string){
		"12_legacy_no_token.yaml": func(t *testing.T, c Configuration, err error, _, after string) {
			require.NoError(t, err, "viper panicked on Profiles[0]")
			assert.Equal(t, currentVersion, c.Version)
			assert.Empty(t, c.Profiles)
			assert.Equal(t, "profiles: []\nqradar:\n    last_execution: 2026-10-11T09:41:12.123456789+02:00\nversion: \"1.0\"\n", after)
		},
		"22_hand_version_unquoted.yaml": func(t *testing.T, c Configuration, err error, input, after string) {
			require.NoError(t, err)
			assert.Equal(t, "1.0", c.Version, "viper read the float 1 as \"1\" and re-ran the migration")
			assert.Equal(t, input, after)
		},
		"23_hand_quoted_time.yaml": func(t *testing.T, c Configuration, err error, input, after string) {
			require.NoError(t, err, "viper could not decode a quoted timestamp")
			assert.Equal(t, []Profile{{Name: "default", Profile: ProfileConfiguration{Host: "https://api.redcarbon.ai", Token: "tok"}}}, c.Profiles)
			assert.Equal(t, input, after)
		},
		"24_hand_profiles_no_version_plus_legacy.yaml": func(t *testing.T, c Configuration, err error, _, _ string) {
			require.NoError(t, err)
			assert.Equal(t, []Profile{
				{Name: "first", Profile: ProfileConfiguration{Host: "https://api.redcarbon.ai", Token: "tok-first"}},
				{Name: "default", Profile: ProfileConfiguration{Token: "tok-legacy"}},
			}, c.Profiles, "viper attached the legacy profile to Profiles[0]")
		},
		"26_hand_weak_types.yaml": func(t *testing.T, c Configuration, err error, input, after string) {
			require.NoError(t, err)
			assert.Equal(t, "0123", c.Profiles[0].Profile.Token, "viper weakly decoded the octal token 0123 as \"83\"")
			assert.Equal(t, input, after)
		},
	}

	files, err := filepath.Glob("testdata/expected/*.viper.json")
	require.NoError(t, err)
	require.NotEmpty(t, files)

	for _, f := range files {
		var want viperResult
		require.NoError(t, json.Unmarshal([]byte(readString(t, f)), &want))

		t.Run(want.File, func(t *testing.T) {
			input := readString(t, filepath.Join("testdata/golden", want.File))
			name := setupFile(t, []byte(input))

			c, err := load()
			after := readString(t, name)

			if check, ok := deviations[want.File]; ok {
				check(t, c, err, input, after)
				return
			}

			require.Empty(t, want.Panic)
			require.NoError(t, err)
			got, err := json.Marshal(c)
			require.NoError(t, err)
			assert.JSONEq(t, string(want.Config), string(got))
			assert.Equal(t, want.FileAfter, after)
		})
	}
}

// Golden files were written by the viper-based implementation running the same
// steps: byte-identical output is what lets an older agent read files written by this one.
func TestWriterMatchesViper(t *testing.T) {
	l := slog.Default()
	dir := t.TempDir()
	require.NoError(t, Init(dir))
	name := filepath.Join(dir, fileName)

	golden := func(t *testing.T, file string) {
		t.Helper()
		assert.Equal(t, readString(t, filepath.Join("testdata/golden", file)), readString(t, name), file)
	}
	addProfile := func(n string, opts ProfileConfiguration) {
		c := LoadConfiguration()
		c.Profiles = append(c.Profiles, Profile{Name: n, Profile: opts})
		c.MustSave()
	}

	golden(t, "01_fresh_init.yaml")

	LoadConfiguration()
	golden(t, "02_fresh_loaded.yaml")

	addProfile("default", ProfileConfiguration{Host: "https://api.redcarbon.ai", Token: "tok-default"})
	golden(t, "03_profile_add.yaml")

	addProfile("Secondary", ProfileConfiguration{Host: "https://eu.redcarbon.ai:443", Token: "0123456789"})
	addProfile("debugger", ProfileConfiguration{Host: "localhost:8080", Token: "tok: with colon"})
	golden(t, "04_profile_add_three.yaml")

	t.Run("re-saving an unchanged profile is byte-stable", func(t *testing.T) {
		OverwriteProfileInConfig(l, LoadProfile("Secondary"))
		golden(t, "04_profile_add_three.yaml")
	})
}

func TestConcurrentOverwrite(t *testing.T) {
	name := setupFile(t, []byte(readString(t, "testdata/golden/04_profile_add_three.yaml")))

	names := []string{"default", "Secondary", "debugger"}

	var wg sync.WaitGroup
	for _, n := range names {
		wg.Go(func() {
			for i := range 20 {
				p := LoadProfile(n)
				p.Profile.Token = fmt.Sprintf("tok-%d", i)
				OverwriteProfileInConfig(slog.Default(), p)
			}
		})
	}
	wg.Wait()

	c := LoadConfiguration()
	require.Len(t, c.Profiles, len(names))
	for i, p := range c.Profiles {
		assert.Equal(t, names[i], p.Name)
		assert.Equal(t, "tok-19", p.Profile.Token, p.Name)
	}
	matches, err := filepath.Glob(name + ".tmp-*")
	require.NoError(t, err)
	assert.Empty(t, matches)
}

// Another agent process (e.g. `redcarbon profile list`) reads the file without
// holding this process' lock: a truncated file would be migrated into an empty one.
func TestReadersNeverObservePartialFile(t *testing.T) {
	want := readString(t, "testdata/golden/04_profile_add_three.yaml")
	name := setupFile(t, []byte(want))

	stop := make(chan struct{})
	var wg sync.WaitGroup
	wg.Go(func() {
		for {
			select {
			case <-stop:
				return
			default:
				OverwriteProfileInConfig(slog.Default(), LoadProfile("default"))
			}
		}
	})

	for range 2000 {
		b, err := os.ReadFile(name)
		require.NoError(t, err)
		if !assert.Equal(t, want, string(b)) {
			break
		}
	}
	close(stop)
	wg.Wait()
}

func TestSaveKeepsSymlinkAndMode(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "real.yaml")
	require.NoError(t, os.WriteFile(target, []byte("profiles: []\nversion: \"1.0\"\n"), 0o600))
	require.NoError(t, os.Symlink(target, filepath.Join(dir, fileName)))
	require.NoError(t, Init(dir))

	c := LoadConfiguration()
	c.Profiles = append(c.Profiles, Profile{Name: "p"})
	c.MustSave()

	fi, err := os.Lstat(filepath.Join(dir, fileName))
	require.NoError(t, err)
	assert.Equal(t, os.ModeSymlink, fi.Mode().Type())
	st, err := os.Stat(target)
	require.NoError(t, err)
	assert.Equal(t, os.FileMode(0o600), st.Mode().Perm())
	assert.Contains(t, readString(t, target), "name: p")
}

func TestSaveInReadOnlyDirectory(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("root ignores directory permissions")
	}
	name := setupFile(t, []byte("profiles: []\nversion: \"1.0\"\n"))
	dir := filepath.Dir(name)
	require.NoError(t, os.Chmod(dir, 0o555))
	t.Cleanup(func() { os.Chmod(dir, 0o755) })

	c := LoadConfiguration()
	c.Profiles = append(c.Profiles, Profile{Name: "p"})
	require.NoError(t, c.save())

	assert.Contains(t, readString(t, name), "name: p")
}

// read lowercases every key, so a camelCase yaml tag would silently stop loading.
func TestEveryFieldSurvivesSaveAndLoad(t *testing.T) {
	setupFile(t, []byte("{}\n"))

	var want Configuration
	fillAllFields(t, reflect.ValueOf(&want).Elem())
	want.Version = currentVersion
	want.MustSave()

	assert.Equal(t, want, LoadConfiguration())
}

func fillAllFields(t *testing.T, v reflect.Value) {
	t.Helper()
	switch {
	case v.Kind() == reflect.String:
		v.SetString("X")
	case v.Kind() == reflect.Slice:
		v.Set(reflect.MakeSlice(v.Type(), 1, 1))
		fillAllFields(t, v.Index(0))
	case v.Kind() == reflect.Struct:
		for i := range v.NumField() {
			fillAllFields(t, v.Field(i))
		}
	default:
		t.Fatalf("fillAllFields does not handle %s", v.Type())
	}
}

// Viper matched keys case-insensitively, so every field must load from an uppercased key.
func TestLoadMatchesKeysCaseInsensitively(t *testing.T) {
	want := Configuration{
		Version: currentVersion,
		Profiles: []Profile{{
			Name:    "P",
			Profile: ProfileConfiguration{Host: "H", Token: "T"},
		}},
	}
	b, err := yaml.Marshal(want)
	require.NoError(t, err)
	setupFile(t, []byte(strings.ToUpper(string(b))))

	got, err := load()
	require.NoError(t, err)
	assert.Equal(t, want, got)
}

// Files written before OPE-2693 still carry per-profile job checkpoints; they must load and be dropped on save.
func TestCheckpointSectionsAreDroppedOnSave(t *testing.T) {
	name := setupFile(t, []byte(readString(t, "testdata/golden/07_time_now_local.yaml")))

	c := LoadConfiguration()
	c.MustSave()

	assert.Equal(t, `profiles:
    - name: default
      profile:
        host: https://api.redcarbon.ai
        token: tok-default
    - name: Secondary
      profile:
        host: https://eu.redcarbon.ai:443
        token: "0123456789"
    - name: debugger
      profile:
        host: localhost:8080
        token: 'tok: with colon'
version: "1.0"
`, readString(t, name))
}
