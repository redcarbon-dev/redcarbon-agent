package config

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"time"
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

func LoadConfiguration() Configuration {
	c := Configuration{}
	err := viper.Unmarshal(&c)
	if err != nil {
		logrus.Fatal("Error while loading the configuration")
	}

	if c.Version != "1.0" {
		c.updateConfig()
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

	logrus.Fatalf("Profile %s not found", name)

	return Profile{}
}

func (c *Configuration) updateConfig() {
	c.Version = "1.0"
	token := viper.GetString("auth.access_token")
	host := viper.GetString("server.host")
	if token != "" {
		c.Profiles = append(c.Profiles, Profile{
			Name: "default",
			Profile: ProfileConfiguration{
				Token: token,
				Host:  host,
			},
		})
	}

	qRadarLastExecution := viper.GetTime("qradar.last_execution")
	if !qRadarLastExecution.IsZero() {
		c.Profiles[0].QRadar = QRadar{
			LastExecution: qRadarLastExecution,
		}
	}

	fortiSIEMLastExecution := viper.GetTime("fortisiem.last_execution")
	if !fortiSIEMLastExecution.IsZero() {
		c.Profiles[0].FortiSIEM = FortiSIEM{
			LastExecution: fortiSIEMLastExecution,
		}
	}

	sentinelONELastExecution := viper.GetTime("sentinel_one.last_execution")
	if !sentinelONELastExecution.IsZero() {
		c.Profiles[0].SentinelONE = SentinelONE{
			LastExecution: sentinelONELastExecution,
		}
	}

	c.MustSave()
}

func (c *Configuration) MustSave() {
	err := c.save()
	if err != nil {
		logrus.Fatal("Error while saving the configuration")
	}
}

func (c *Configuration) save() error {
	viper.Set("version", c.Version)
	viper.Set("profiles", c.Profiles)

	return viper.WriteConfig()
}

func OverwriteProfileInConfig(l *logrus.Entry, profile Profile) {
	c := LoadConfiguration()

	for i, p := range c.Profiles {
		if p.Name == profile.Name {
			c.Profiles[i] = profile
		}
	}

	err := c.save()
	if err != nil {
		l.WithError(err).Error("failed to save profile")
	}
}
