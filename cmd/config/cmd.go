package config

import (
	"crypto/tls"
	"os"
	"path"
	agentsPublicApiV1 "pkg.redcarbon.ai/proto/redcarbon/public_apis/agents/api/v1"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"

	"pkg.redcarbon.ai/internal/auth"
	"pkg.redcarbon.ai/internal/build"
)

type ConfigOptions struct {
	RefreshToken string
	Host         string
	Insecure     bool
}

func NewConfigCmd() *cobra.Command {
	opts := &ConfigOptions{}

	cmd := &cobra.Command{
		Use:   "config",
		Short: "Configure the RedCarbon agent",
		Run: func(cmd *cobra.Command, args []string) {
			run(cmd, args, opts)
		},
	}

	cmd.Flags().StringVarP(&opts.RefreshToken, "token", "t", "", "The token used to execute the login")
	cmd.Flags().StringVarP(&opts.Host, "server", "s", build.DefaultHost, "The Server used to execute the login")
	cmd.Flags().BoolVarP(&opts.Insecure, "insecure", "i", false, "connection insecure")

	return cmd
}

func run(cmd *cobra.Command, args []string, opts *ConfigOptions) {
	viper.Set("server.host", opts.Host)
	viper.Set("server.insecure", opts.Insecure)

	agentsCli := mustCreateAgentCli()
	defer agentsCli.Close()

	confDir, err := os.UserConfigDir()
	if err != nil {
		logrus.Fatalf("can't extract the user config directory for error %v", err)
	}

	redcarbonConfDir := path.Join(confDir, "redcarbon")

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(redcarbonConfDir)

	client := agentsPublicApiV1.NewAgentsPublicApiV1SrvClient(agentsCli)
	a := auth.NewAuthService(client, path.Join(confDir, "redcarbon", "config.yaml"))

	err = a.RefreshToken(opts.RefreshToken)
	if err != nil {
		logrus.Fatalf("Can't refresh the token %v", err)
	}

	logrus.Println("Agent successfully configured!")
}

func mustCreateAgentCli() *grpc.ClientConn {
	host := viper.GetString("server.host")

	var creds credentials.TransportCredentials

	logrus.Infof("Connecting to the server... %s\n", host)

	if viper.GetBool("server.insecure") {
		creds = insecure.NewCredentials()
	} else {
		creds = credentials.NewTLS(&tls.Config{})
	}

	agentsCli, err := grpc.Dial(host, grpc.WithTransportCredentials(creds))
	if err != nil {
		logrus.Fatalf("Cannot create source connection: %v", err)
	}

	return agentsCli
}
