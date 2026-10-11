package config

import (
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"go.yaml.in/yaml/v3"
)

const fileName = ".config.yaml"

const currentVersion = "1.0"

var (
	mu         sync.Mutex
	configPath = fileName
)

type Configuration struct {
	Version  string    `yaml:"version"`
	Profiles []Profile `yaml:"profiles"`
}

type Profile struct {
	Name        string               `yaml:"name"`
	Profile     ProfileConfiguration `yaml:"profile"`
	SentinelONE SentinelONE          `yaml:"sentinelone,omitempty"`
	FortiSIEM   FortiSIEM            `yaml:"fortisiem,omitempty"`
	QRadar      QRadar               `yaml:"qradar,omitempty"`
	Debug       Debug                `yaml:"debug,omitempty"`
}

type ProfileConfiguration struct {
	Host  string `yaml:"host"`
	Token string `yaml:"token"`
}

type Debug struct {
	Active        bool      `yaml:"active"`
	LastExecution time.Time `yaml:"lastExecution"`
}

type SentinelONE struct {
	LastExecution time.Time `yaml:"lastExecution"`
}

type FortiSIEM struct {
	LastExecution time.Time `yaml:"lastExecution"`
}

type QRadar struct {
	LastExecution time.Time `yaml:"lastExecution"`
}

// Init points the package at the configuration file in dir, creating an empty one if it does not exist.
func Init(dir string) error {
	mu.Lock()
	defer mu.Unlock()

	configPath = filepath.Join(dir, fileName)

	f, err := os.OpenFile(configPath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o644)
	if errors.Is(err, fs.ErrExist) {
		return nil
	}
	if err != nil {
		return err
	}

	_, err = f.WriteString("{}\n")

	return errors.Join(err, f.Close())
}

func LoadConfiguration() Configuration {
	mu.Lock()
	defer mu.Unlock()

	c, err := load()
	if err != nil {
		slog.Error("Error while loading the configuration", "error", err)
		os.Exit(1)
	}

	return c
}

func LoadProfile(name string) Profile {
	c := LoadConfiguration()

	for _, p := range c.Profiles {
		if p.Name == name {
			return p
		}
	}

	slog.Error("Profile not found", "profile", name)
	os.Exit(1)

	return Profile{}
}

func (c *Configuration) MustSave() {
	mu.Lock()
	defer mu.Unlock()

	if err := c.save(); err != nil {
		slog.Error("Error while saving the configuration", "error", err)
		os.Exit(1)
	}
}

func OverwriteProfileInConfig(l *slog.Logger, profile Profile) {
	mu.Lock()
	defer mu.Unlock()

	c, err := load()
	if err != nil {
		l.Error("failed to load the configuration", "error", err)
		return
	}

	for i, p := range c.Profiles {
		if p.Name == profile.Name {
			c.Profiles[i] = profile
		}
	}

	if err := c.save(); err != nil {
		l.Error("failed to save profile", "error", err)
	}
}

func load() (Configuration, error) {
	var c Configuration

	doc, node, err := read()
	if err != nil {
		return c, err
	}

	if err := node.Decode(&c); err != nil {
		return c, fmt.Errorf("decoding %s: %w", configPath, err)
	}

	if c.Version != currentVersion {
		c.migrateLegacy(node)
		if err := c.write(doc); err != nil {
			return c, err
		}
	}

	return c, nil
}

// read returns the file both as a generic document, whose unknown and legacy
// keys are carried over on save, and as a node to decode the configuration from.
// Keys are lowercased because viper matched them case-insensitively and wrote
// them back lowercased.
func read() (map[string]any, *yaml.Node, error) {
	b, err := os.ReadFile(configPath)
	if err != nil {
		return nil, nil, err
	}

	var node yaml.Node
	if err := yaml.Unmarshal(b, &node); err != nil {
		return nil, nil, fmt.Errorf("parsing %s: %w", configPath, err)
	}

	renameKeys(&node, strings.ToLower)

	var doc map[string]any
	if err := node.Decode(&doc); err != nil {
		return nil, nil, fmt.Errorf("parsing %s: %w", configPath, err)
	}
	if doc == nil {
		doc = map[string]any{}
	}

	// lastExecution is the only yaml tag that is not lowercase;
	// TestEveryFieldSurvivesSaveAndLoad fails if another one is added.
	renameKeys(&node, func(k string) string {
		if k == "lastexecution" {
			return "lastExecution"
		}
		return k
	})

	return doc, &node, nil
}

func renameKeys(n *yaml.Node, rename func(string) string) {
	if n.Kind == yaml.MappingNode {
		for i := 0; i < len(n.Content); i += 2 {
			n.Content[i].Value = rename(n.Content[i].Value)
		}
	}

	for _, c := range n.Content {
		renameKeys(c, rename)
	}
}

func (c *Configuration) save() error {
	doc, _, err := read()
	if err != nil {
		return err
	}

	return c.write(doc)
}

func (c *Configuration) write(doc map[string]any) error {
	doc["version"] = c.Version
	doc["profiles"] = c.Profiles

	b, err := yaml.Marshal(doc)
	if err != nil {
		return err
	}

	return writeFile(configPath, b)
}

// writeFile replaces the file atomically (temp file + rename) so that a crash or a
// concurrent reader never observes a truncated configuration, which would be
// migrated into an empty one and lose every profile. If the directory is not
// writable it falls back to rewriting the file in place, which is NOT atomic.
func writeFile(name string, data []byte) error {
	if resolved, err := filepath.EvalSymlinks(name); err == nil {
		name = resolved
	}

	perm := fs.FileMode(0o644)
	if st, err := os.Stat(name); err == nil {
		perm = st.Mode().Perm()
	}

	tmp, err := os.CreateTemp(filepath.Dir(name), filepath.Base(name)+".tmp-*")
	if errors.Is(err, fs.ErrPermission) {
		// The directory may be read-only for the agent user while the file itself is
		// writable: a non-atomic write beats failing to persist LastExecution at all.
		return os.WriteFile(name, data, perm)
	}
	if err != nil {
		return err
	}
	defer os.Remove(tmp.Name())

	_, err = tmp.Write(data)
	if err == nil {
		err = tmp.Chmod(perm)
	}
	if err == nil {
		err = tmp.Sync()
	}
	if err = errors.Join(err, tmp.Close()); err != nil {
		return err
	}

	return os.Rename(tmp.Name(), name)
}

type legacyConfiguration struct {
	Auth struct {
		AccessToken string `yaml:"access_token"`
	}
	Server struct {
		Host string
	}
	QRadar      legacyService
	FortiSIEM   legacyService
	SentinelOne legacyService `yaml:"sentinel_one"`
}

type legacyService struct {
	LastExecution time.Time `yaml:"last_execution"`
}

// migrateLegacy converts the flat format written by agents before profiles were introduced.
func (c *Configuration) migrateLegacy(node *yaml.Node) {
	c.Version = currentVersion

	var legacy legacyConfiguration
	// Best effort, as viper's getters were: a malformed legacy key must not prevent loading.
	_ = node.Decode(&legacy)

	if legacy.Auth.AccessToken == "" {
		return
	}

	c.Profiles = append(c.Profiles, Profile{
		Name:        "default",
		Profile:     ProfileConfiguration{Host: legacy.Server.Host, Token: legacy.Auth.AccessToken},
		QRadar:      QRadar{LastExecution: legacy.QRadar.LastExecution},
		FortiSIEM:   FortiSIEM{LastExecution: legacy.FortiSIEM.LastExecution},
		SentinelONE: SentinelONE{LastExecution: legacy.SentinelOne.LastExecution},
	})
}
