package profile

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"pkg.redcarbon.ai/internal/config"
)

func newProfileRmCmd() *cobra.Command {
	opts := config.ProfileConfiguration{}

	cmd := &cobra.Command{
		Use:   "rm [profile]",
		Short: "Remove an existent profile",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return fmt.Errorf("requires a profile name")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			runRm(cmd, args, opts)
		},
	}

	return cmd
}

func runRm(cmd *cobra.Command, args []string, opts config.ProfileConfiguration) {
	conf := config.LoadConfiguration()

	profileName := args[0]

	found := -1

	for i, p := range conf.Profiles {
		if p.Name == profileName {
			found = i
			break
		}
	}

	if found == -1 {
		logrus.Fatalf("Profile %s not found", profileName)
	}

	conf.Profiles = append(conf.Profiles[:found], conf.Profiles[found+1:]...)

	conf.MustSave()

	logrus.Println("Profile successfully removed!")
}
