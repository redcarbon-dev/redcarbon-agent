package main

import (
	"os"
	"pkg.redcarbon.ai/internal/build"

	"github.com/spf13/cobra"

	"pkg.redcarbon.ai/cmd/config"
	"pkg.redcarbon.ai/cmd/start"
)

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
