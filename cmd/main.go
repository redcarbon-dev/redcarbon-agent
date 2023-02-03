package main

import (
	"os"
	"path"

	"pkg.redcarbon.ai/internal/build"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"pkg.redcarbon.ai/cmd/config"
	"pkg.redcarbon.ai/cmd/start"
)

func init() {
	confDir, err := os.UserConfigDir()
	if err != nil {
		logrus.Fatalf("can't extract the user config directory for error %v", err)
	}

	redcarbonConfDir := path.Join(confDir, "redcarbon")

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(redcarbonConfDir)

	err = viper.ReadInConfig()
	if err != nil {
		logrus.Fatalf("can't read the configuration %v", err)
	}

	viper.SetDefault("server.host", "localhost:50051")
	viper.SetDefault("server.insecure", true)
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
