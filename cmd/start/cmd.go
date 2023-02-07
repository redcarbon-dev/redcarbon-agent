package start

import (
	"crypto/tls"
	"os"
	"path"

	"github.com/jasonlvhit/gocron"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"

	"pkg.redcarbon.ai/internal/auth"
	"pkg.redcarbon.ai/internal/routines"
	agentsExternalApiV1 "pkg.redcarbon.ai/proto/redcarbon/external_api/agents/api/v1"
)

const (
	hzRoutineInterval      = 5
	refreshRoutineInterval = 30
	configRoutineInterval  = 10
)

func NewStartCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "start",
		Short: "Start the RedCarbon agent",
		Run:   run,
	}
}

func run(cmd *cobra.Command, args []string) {
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

	err = viper.ReadInConfig()
	if err != nil {
		logrus.Fatalf("can't read the configuration %v", err)
	}

	client := agentsExternalApiV1.NewAgentsExternalV1SrvClient(agentsCli)
	a := auth.NewAuthService(client, path.Join(confDir, "redcarbon", "config.yaml"))
	r := routines.NewRoutineJobs(client, a)

	gocron.Every(hzRoutineInterval).Seconds().From(gocron.NextTick()).Do(r.HZRoutine)
	gocron.Every(refreshRoutineInterval).Minutes().From(gocron.NextTick()).Do(r.Refresh)
	gocron.Every(configRoutineInterval).Minutes().From(gocron.NextTick()).Do(r.ConfigRoutine)

	<-gocron.Start()
}

func mustCreateAgentCli() *grpc.ClientConn {
	host := viper.GetString("server.host")

	var creds credentials.TransportCredentials

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
