package config

import (
	"crypto/tls"
	"os"
	"path"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"

	"pkg.redcarbon.ai/internal/auth"
	"pkg.redcarbon.ai/internal/build"
	agentsExternalApiV1 "pkg.redcarbon.ai/proto/redcarbon/external_api/agents/api/v1"
)

type ConfigOptions struct {
	RefreshToken string
	Host         string
	Insercure    bool
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
	cmd.Flags().BoolVarP(&opts.Insercure, "insecure", "i", false, "connection insecure")

	return cmd
}

func run(cmd *cobra.Command, args []string, opts *ConfigOptions) {
	viper.Set("server.host", opts.Host)
	viper.Set("server.insecure", opts.Insercure)

	agentsCli := mustCreateAgentCli()
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
