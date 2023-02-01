package start

import (
	"github.com/jasonlvhit/gocron"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"os"
	"path"
	"pkg.redcarbon.ai/internal/auth"

	"pkg.redcarbon.ai/internal/routines"
	agentsExternalApiV1 "pkg.redcarbon.ai/proto/redcarbon/external_api/agents/api/v1"
)

func NewStartCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "start",
		Short: "Start the RedCarbon agent",
		Run:   run,
	}
}

func run(cmd *cobra.Command, args []string) {
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

	gocron.Every(5).Seconds().From(gocron.NextTick()).Do(r.HZRoutine)
	gocron.Every(30).Seconds().From(gocron.NextTick()).Do(r.Refresh)

	<-gocron.Start()
}
