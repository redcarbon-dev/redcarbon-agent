package config

import (
	"encoding/json"
	"log/slog"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"

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
// wrote back to, each golden input. Panic is set where viper crashed or exited.
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
			assert.True(t, c.Profiles[0].QRadar.LastExecution.Equal(time.Date(2026, 10, 11, 9, 41, 12, 0, time.UTC)))
			assert.Equal(t, input, after)
		},
		"24_hand_profiles_no_version_plus_legacy.yaml": func(t *testing.T, c Configuration, err error, _, _ string) {
			require.NoError(t, err)
			require.Len(t, c.Profiles, 2)
			assert.True(t, c.Profiles[0].QRadar.LastExecution.IsZero(), "viper attached the legacy timestamp to Profiles[0]")
			assert.Equal(t, "default", c.Profiles[1].Name)
			assert.False(t, c.Profiles[1].QRadar.LastExecution.IsZero())
		},
		"26_hand_weak_types.yaml": func(t *testing.T, _ Configuration, err error, input, after string) {
			assert.Error(t, err, "viper weakly decoded active: \"true\" and the octal token 0123 as \"83\"")
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

var (
	cest    = time.FixedZone("CEST", 2*60*60)
	tQRadar = time.Date(2026, 10, 11, 9, 41, 12, 123456789, cest)
	tS1     = time.Date(2026, 10, 11, 7, 0, 0, 0, time.UTC)
	tForti  = time.Date(2026, 10, 10, 23, 59, 59, 999000000, time.UTC)
	tDebug  = time.Date(2026, 1, 15, 12, 30, 0, 500, time.FixedZone("EST", -5*60*60))
)

// Golden files were written by the viper-based implementation running the same
// steps: byte-identical output is what lets an older agent read files written by this one.
func TestWriterMatchesViper(t *testing.T) {
	l := slog.Default()
	dir := t.TempDir()
	require.NoError(t, Init(dir))
	name := filepath.Join(dir, fileName)

	golden := func(file string) {
		t.Helper()
		assert.Equal(t, readString(t, filepath.Join("testdata/golden", file)), readString(t, name), file)
	}
	addProfile := func(n string, opts ProfileConfiguration) {
		c := LoadConfiguration()
		c.Profiles = append(c.Profiles, Profile{Name: n, Profile: opts})
		c.MustSave()
	}

	golden("01_fresh_init.yaml")

	LoadConfiguration()
	golden("02_fresh_loaded.yaml")

	addProfile("default", ProfileConfiguration{Host: "https://api.redcarbon.ai", Token: "tok-default"})
	golden("03_profile_add.yaml")

	addProfile("Secondary", ProfileConfiguration{Host: "https://eu.redcarbon.ai:443", Token: "0123456789"})
	addProfile("debugger", ProfileConfiguration{Host: "localhost:8080", Token: "tok: with colon"})
	golden("04_profile_add_three.yaml")

	p := LoadProfile("default")
	p.QRadar.LastExecution = tQRadar
	OverwriteProfileInConfig(l, p)
	p = LoadProfile("Secondary")
	p.SentinelONE.LastExecution = tS1
	p.FortiSIEM.LastExecution = tForti
	OverwriteProfileInConfig(l, p)
	p = LoadProfile("debugger")
	p.Debug.Active = true
	OverwriteProfileInConfig(l, p)
	golden("05_services_updated.yaml")

	p = LoadProfile("debugger")
	p.Debug.LastExecution = tDebug
	OverwriteProfileInConfig(l, p)
	golden("06_debug_updated.yaml")
}

func TestConcurrentOverwrite(t *testing.T) {
	name := setupFile(t, []byte(readString(t, "testdata/golden/04_profile_add_three.yaml")))

	names := []string{"default", "Secondary", "debugger"}

	var wg sync.WaitGroup
	for _, n := range names {
		wg.Go(func() {
			for i := range 20 {
				p := LoadProfile(n)
				p.QRadar.LastExecution = tS1.Add(time.Duration(i) * time.Second)
				OverwriteProfileInConfig(slog.Default(), p)
			}
		})
	}
	wg.Wait()

	c := LoadConfiguration()
	require.Len(t, c.Profiles, len(names))
	for i, p := range c.Profiles {
		assert.Equal(t, names[i], p.Name)
		assert.True(t, p.QRadar.LastExecution.Equal(tS1.Add(19*time.Second)), p.Name)
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

// read lowercases every key and then restores the camelCase yaml tags it knows
// about, so a new camelCase tag it does not know would silently stop loading.
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
	case v.Type() == reflect.TypeFor[time.Time]():
		v.Set(reflect.ValueOf(tS1))
	case v.Kind() == reflect.String:
		v.SetString("X")
	case v.Kind() == reflect.Bool:
		v.SetBool(true)
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
			Name:        "P",
			Profile:     ProfileConfiguration{Host: "H", Token: "T"},
			SentinelONE: SentinelONE{LastExecution: tS1},
			FortiSIEM:   FortiSIEM{LastExecution: tForti},
			QRadar:      QRadar{LastExecution: tS1},
			Debug:       Debug{Active: true, LastExecution: tForti},
		}},
	}
	b, err := yaml.Marshal(want)
	require.NoError(t, err)
	setupFile(t, []byte(strings.ToUpper(string(b))))

	got, err := load()
	require.NoError(t, err)
	assert.Equal(t, want, got)
}
