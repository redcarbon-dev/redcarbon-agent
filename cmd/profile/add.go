package profile

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"pkg.redcarbon.ai/internal/build"
	"pkg.redcarbon.ai/internal/config"
)

type ConfigOptions struct {
	Token string
	Host  string
}

func newProfileAddCmd() *cobra.Command {
	opts := config.ProfileConfiguration{}

	cmd := &cobra.Command{
		Use:   "add [profile]",
		Short: "Add a new profile",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return fmt.Errorf("requires a profile name")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			runAdd(cmd, args, opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Token, "token", "t", "", "The token used to execute the login")
	cmd.Flags().StringVarP(&opts.Host, "server", "s", build.DefaultHost, "The Server used to execute the login")

	cmd.MarkFlagRequired("token")

	return cmd
}

func runAdd(cmd *cobra.Command, args []string, opts config.ProfileConfiguration) {
	conf := config.LoadConfiguration()

	profileName := args[0]

	for _, p := range conf.Profiles {
		if p.Name == profileName {
			logrus.Fatalf("Profile %s already exists", profileName)
		}
	}

	profile := config.Profile{
		Name:    profileName,
		Profile: opts,
	}

	conf.Profiles = append(conf.Profiles, profile)

	conf.MustSave()

	logrus.Println("Profile successfully added!")
}
