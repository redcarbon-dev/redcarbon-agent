package start

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"pkg.redcarbon.ai/proto/redcarbon/agents_public/v1/agents_publicv1connect"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/google/go-github/v50/github"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"pkg.redcarbon.ai/internal/build"
	"pkg.redcarbon.ai/internal/routines"
)

const (
	hzRoutineInterval     = 5
	configRoutineInterval = 10
	updateRoutineInterval = 1
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

	if viper.GetString("auth.access_token") == "" {
		logrus.Fatal("No access token found. Please run `redcarbon config` to configure the agent")
	}

	ctx, cancelFn := context.WithCancel(context.Background())
	defer cancelFn()

	host := viper.GetString("server.host")

	client := agents_publicv1connect.NewAgentsPublicAPIsV1SrvClient(http.DefaultClient, host)
	gh := github.NewClient(nil)
	done := make(chan bool)
	r := routines.NewRoutineJobs(client, gh, done)

	s := gocron.NewScheduler(time.UTC)

	s.Every(updateRoutineInterval).Day().StartImmediately().SingletonMode().Do(r.UpdateRoutine, ctx)
	s.Every(hzRoutineInterval).Seconds().StartImmediately().Do(r.HZRoutine, ctx)
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
		s.Stop()
		logrus.Info("RedCarbon Agent stopped")
		os.Exit(0)
	case <-done:
		s.Stop()
		logrus.Info("RedCarbon Agent stopped due update")
		os.Exit(3)
	}
}
