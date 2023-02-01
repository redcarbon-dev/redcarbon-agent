package config

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"os"
	"path"

	"pkg.redcarbon.ai/internal/auth"
	agentsExternalApiV1 "pkg.redcarbon.ai/proto/redcarbon/external_api/agents/api/v1"
)

type LoginOptions struct {
	RefreshToken string
}

func NewConfigCmd() *cobra.Command {
	opts := &LoginOptions{}

	cmd := &cobra.Command{
		Use:   "config",
		Short: "Configure the RedCarbon agent",
		Run: func(cmd *cobra.Command, args []string) {
			run(cmd, args, opts)
		},
	}

	cmd.Flags().StringVarP(&opts.RefreshToken, "token", "t", "", "The token used to execute the login")

	return cmd
}

func run(cmd *cobra.Command, args []string, opts *LoginOptions) {
	agentsCli, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logrus.Fatalf("Cannot create source connection: %v", err)
	}
	defer agentsCli.Close()

	confDir, err := os.UserConfigDir()
	if err != nil {
		logrus.Fatalf("can't extract the user config directory for error %v", err)
	}

	redcarbonConfDir := path.Join(confDir, "redcarbon")

	err = os.MkdirAll(redcarbonConfDir, 0755)
	if err != nil {
		logrus.Fatalf("can't create redcarbon config directory for error %v", err)
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(redcarbonConfDir)

	client := agentsExternalApiV1.NewAgentsExternalV1SrvClient(agentsCli)
	a := auth.NewAuthService(client, path.Join(confDir, "redcarbon", "config.yaml"))

	err = a.RefreshToken(opts.RefreshToken)
	if err != nil {
		logrus.Fatalf("Can't refresh the token %v", err)
	}

	logrus.Println("Agent successfully configured!")
}
