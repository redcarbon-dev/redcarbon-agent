package profile

import "github.com/spf13/cobra"

func NewProfileCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "profile",
		Short: "Handle profiles for the RedCarbon agent",
	}

	cmd.AddCommand(newProfileAddCmd())
	cmd.AddCommand(newProfileListCmd())
	cmd.AddCommand(newProfileRmCmd())

	return cmd
}
