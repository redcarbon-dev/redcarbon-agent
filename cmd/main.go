package main

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"
	"pkg.redcarbon.ai/cmd/profile"
	"pkg.redcarbon.ai/internal/cli"
	"pkg.redcarbon.ai/internal/config"
	"pkg.redcarbon.ai/internal/github"
	"pkg.redcarbon.ai/internal/routines"
	"pkg.redcarbon.ai/internal/scheduler"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"pkg.redcarbon.ai/internal/build"
)

const (
	hzRoutineInterval     = 5 * time.Second
	configRoutineInterval = 10 * time.Minute
	updateRoutineInterval = 24 * time.Hour
	debugRoutineInterval  = 2 * time.Second
	proxyRoutineInterval  = 500 * time.Millisecond

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

	rootCmd.AddCommand(profile.NewProfileCmd())

	rootCmd.Version = build.Version

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, args []string) {
	logrus.Infof("Starting RedCarbon Agent v%s on %s", build.Version, build.Architecture)

	conf := config.LoadConfiguration()

	if len(conf.Profiles) == 0 {
		logrus.Fatal("No profiles found, please add one by running `redcarbon profile add`")
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT)
	defer cancel()

	g, ctx := errgroup.WithContext(ctx)

	clientFactory := cli.NewClientFactory()
	gh := github.NewClient(nil, github.DefaultBaseURL)
	done := make(chan bool)

	var jobs []scheduler.Job
	for _, prof := range conf.Profiles {
		configRoutine := configRoutineInterval
		if prof.Debug.Active {
			configRoutine = debugRoutineInterval
		}

		r := routines.NewRoutineJobs(prof, clientFactory, gh, done)

		jobs = append(jobs,
			scheduler.Job{Name: prof.Name + "/update", Interval: updateRoutineInterval, Singleton: true, Run: r.UpdateRoutine},
			scheduler.Job{Name: prof.Name + "/hz", Interval: hzRoutineInterval, Run: r.HZRoutine},
			scheduler.Job{Name: prof.Name + "/config", Interval: configRoutine, Singleton: true, Run: r.ConfigRoutine},
			scheduler.Job{Name: prof.Name + "/proxy", Interval: proxyRoutineInterval, Run: r.ProxyRoutine},
		)
	}

	s, err := scheduler.Start(ctx, jobs...)
	if err != nil {
		logrus.WithError(err).Fatal("error while scheduling the routines")
	}

	g.Go(func() error {
		<-ctx.Done()
		s.Stop()
		logrus.Info("RedCarbon Agent stopped")
		os.Exit(0)
		return nil
	})

	g.Go(func() error {
		<-done
		s.Stop()
		logrus.Info("RedCarbon Agent stopped due update")
		os.Exit(updateErrorCode)
		return nil
	})

	if err := g.Wait(); err != nil {
		logrus.WithError(err).Fatal("error while running the agent")
	}
}
