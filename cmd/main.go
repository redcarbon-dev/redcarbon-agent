package main

import (
	"context"
	"errors"
	"github.com/go-co-op/gocron"
	"github.com/google/go-github/v50/github"
	"net/http"
	"os"
	"os/signal"
	"pkg.redcarbon.ai/cmd/configure"
	"pkg.redcarbon.ai/internal/routines"
	"pkg.redcarbon.ai/proto/redcarbon/agents_public/v1/agents_publicv1connect"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"pkg.redcarbon.ai/internal/build"
)

const (
	hzRoutineInterval     = 5
	configRoutineInterval = 10
	updateRoutineInterval = 1

	updateErrorCode = 3
)

func init() {
	confDir, err := os.Getwd()
	if err != nil {
		logrus.Fatalf("can't extract the user working directory for error %v", err)
	}

	viper.SetConfigName(".config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(confDir)

	if err = viper.ReadInConfig(); err != nil {
		var e viper.ConfigFileNotFoundError
		if !errors.As(err, &e) {
			logrus.WithError(err).Fatal("error while reading the configuration")
		}

		if err = viper.SafeWriteConfig(); err != nil {
			logrus.WithError(err).Fatal("error while writing the configuration")
		}
	}
}

func main() {
	rootCmd := &cobra.Command{
		Use:   "redcarbon",
		Short: "RedCarbon Agent",
		Run:   run,
	}

	rootCmd.AddCommand(configure.NewConfigureCmd())

	rootCmd.Version = build.Version

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
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
		os.Exit(updateErrorCode)
	}
}
