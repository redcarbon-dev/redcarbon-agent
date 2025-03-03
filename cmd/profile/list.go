package profile

import (
	"fmt"
	"github.com/spf13/cobra"
	"pkg.redcarbon.ai/internal/config"
)

func newProfileListCmd() *cobra.Command {
	opts := config.ProfileConfiguration{}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all profiles",
		Run: func(cmd *cobra.Command, args []string) {
			runList(cmd, args, opts)
		},
	}

	return cmd
}

func runList(cmd *cobra.Command, args []string, opts config.ProfileConfiguration) {
	conf := config.LoadConfiguration()

	for _, p := range conf.Profiles {
		fmt.Println(p.Name)
	}
}
