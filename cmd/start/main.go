package start

import (
	"github.com/jasonlvhit/gocron"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

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

	client := agentsExternalApiV1.NewAgentsExternalV1SrvClient(agentsCli)
	r := routines.NewRoutineJobs(client)

	gocron.Every(5).Seconds().From(gocron.NextTick()).Do(r.HZRoutine)

	<-gocron.Start()
}
