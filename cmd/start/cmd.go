package start

import (
	"context"
	"crypto/tls"
	"os"
	"os/signal"
	"path"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/google/go-github/v50/github"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"

	"pkg.redcarbon.ai/internal/auth"
	"pkg.redcarbon.ai/internal/build"
	"pkg.redcarbon.ai/internal/routines"
	agentsExternalApiV1 "pkg.redcarbon.ai/proto/redcarbon/external_api/agents/api/v1"
)

const (
	hzRoutineInterval      = 5
	refreshRoutineInterval = 30
	configRoutineInterval  = 10
	updateRoutineInterval  = 1
)

func NewStartCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "start",
		Short: "Start the RedCarbon agent",
		Run:   run,
	}
}

func run(cmd *cobra.Command, args []string) {
	logrus.Infof("Starting RedCarbon Agent v%s on %s", build.Version, build.Architecture)

	ctx, cancelFn := context.WithCancel(context.Background())
	defer cancelFn()

	agentsCli := mustCreateAgentCli()
	defer agentsCli.Close()

	confDir, err := os.UserConfigDir()
	if err != nil {
		logrus.Fatalf("can't extract the user config directory for error %v", err)
	}

	client := agentsExternalApiV1.NewAgentsExternalV1SrvClient(agentsCli)
	a := auth.NewAuthService(client, path.Join(confDir, "redcarbon", "config.yaml"))
	gh := github.NewClient(nil)
	done := make(chan bool)
	r := routines.NewRoutineJobs(client, a, gh, done)

	s := gocron.NewScheduler(time.UTC)

	s.Every(updateRoutineInterval).Day().StartImmediately().SingletonMode().Do(r.UpdateRoutine, ctx)
	s.Every(hzRoutineInterval).Seconds().StartImmediately().Do(r.HZRoutine, ctx)
	s.Every(refreshRoutineInterval).Minutes().StartImmediately().SingletonMode().Do(r.Refresh, ctx)
	s.Every(configRoutineInterval).Minutes().StartImmediately().SingletonMode().Do(r.ConfigRoutine, ctx)

	s.StartAsync()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		<-c
		cancelFn()
	}()

	select {
	case <-ctx.Done():
	case <-done:
	}

	s.Stop()

	logrus.Info("RedCarbon Agent stopped")
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
