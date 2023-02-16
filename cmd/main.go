package main

import (
	"errors"
	"os"
	"path"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"pkg.redcarbon.ai/cmd/config"
	"pkg.redcarbon.ai/cmd/start"
	"pkg.redcarbon.ai/internal/build"
)

const confDirPermission = 0o755

func init() {
	confDir, err := os.UserConfigDir()
	if err != nil {
		logrus.Fatalf("can't extract the user config directory for error %v", err)
	}

	redcarbonConfDir := path.Join(confDir, "redcarbon")

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(redcarbonConfDir)

	viper.SetDefault("server.host", build.DefaultHost)
	viper.SetDefault("server.insecure", true)

	if _, err := os.Stat(path.Join(redcarbonConfDir, "config.yaml")); errors.Is(err, os.ErrNotExist) {
		err = os.MkdirAll(redcarbonConfDir, confDirPermission)
		if err != nil {
			logrus.Fatalf("can't create redcarbon config directory for error %v", err)
		}

		err = viper.WriteConfigAs(path.Join(redcarbonConfDir, "config.yaml"))
		if err != nil {
			logrus.Fatalf("can't create redcarbon config file for error %v", err)
		}
	}

	err = viper.ReadInConfig()
	if err != nil {
		logrus.Fatalf("can't read the configuration %v", err)
	}
}

func main() {
	rootCmd := &cobra.Command{
		Use:   "redcarbon",
		Short: "RedCarbon Agent",
	}

	rootCmd.AddCommand(config.NewConfigCmd())
	rootCmd.AddCommand(start.NewStartCmd())

	rootCmd.Version = build.Version

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
