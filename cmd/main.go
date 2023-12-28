package main

import (
	"errors"
	"os"
	"pkg.redcarbon.ai/cmd/configure"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"pkg.redcarbon.ai/cmd/start"
	"pkg.redcarbon.ai/internal/build"
)

const confDirPermission = 0o755

func init() {
	confDir, err := os.Getwd()
	if err != nil {
		logrus.Fatalf("can't extract the user working directory for error %v", err)
	}

	viper.SetConfigName(".config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(confDir)

	if err = viper.ReadInConfig(); err != nil {
		var e viper.ConfigFileNotFoundError
		if !errors.As(err, &e) {
			logrus.WithError(err).Fatal("error while reading the configuration")
		}

		if err = viper.SafeWriteConfig(); err != nil {
			logrus.WithError(err).Fatal("error while writing the configuration")
		}
	}
}

func main() {
	rootCmd := &cobra.Command{
		Use:   "redcarbon",
		Short: "RedCarbon Agent",
	}

	rootCmd.AddCommand(configure.NewConfigureCmd())
	rootCmd.AddCommand(start.NewStartCmd())

	rootCmd.Version = build.Version

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
