package config

import "github.com/spf13/cobra"

func NewConfigCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "config",
		Short: "Configure the RedCarbon agent",
		Run:   run,
	}
}

func run(cmd *cobra.Command, args []string) {

}
