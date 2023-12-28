package configure

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"pkg.redcarbon.ai/internal/build"
)

type ConfigOptions struct {
	Token string
	Host  string
}

func NewConfigureCmd() *cobra.Command {
	opts := &ConfigOptions{}

	cmd := &cobra.Command{
		Use:   "configure",
		Short: "Configure the RedCarbon agent",
		Run: func(cmd *cobra.Command, args []string) {
			run(cmd, args, opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Token, "token", "t", "", "The token used to execute the login")
	cmd.Flags().StringVarP(&opts.Host, "server", "s", build.DefaultHost, "The Server used to execute the login")

	return cmd
}

func run(cmd *cobra.Command, args []string, opts *ConfigOptions) {
	viper.Set("server.host", opts.Host)
	viper.Set("auth.access_token", opts.Token)

	if err := viper.WriteConfig(); err != nil {
		logrus.Fatalf("Error while writing the configuration: %v", err)
	}

	logrus.Println("Agent successfully configured!")
}
