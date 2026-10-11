package main

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/google/go-github/v92/github"
	"golang.org/x/sync/errgroup"
	"pkg.redcarbon.ai/cmd/profile"
	"pkg.redcarbon.ai/internal/cli"
	"pkg.redcarbon.ai/internal/config"
	"pkg.redcarbon.ai/internal/routines"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"pkg.redcarbon.ai/internal/build"
)

const (
	hzRoutineInterval     = "5s"
	configRoutineInterval = "10m"
	updateRoutineInterval = "24h"
	debugRoutineInterval  = "2s"
	proxyRoutineInterval  = "500ms"

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

	s := gocron.NewScheduler(time.UTC)

	clientFactory := cli.NewClientFactory()
	gh, err := github.NewClient()
	if err != nil {
		logrus.WithError(err).Fatal("error while creating the GitHub client")
	}
	done := make(chan bool)

	// Jobs must be registered before starting the scheduler: registering them
	// concurrently with StartAsync lets the scheduler pick up a half-built job
	// and silently drop it.
	for _, prof := range conf.Profiles {
		configRoutine := configRoutineInterval
		if prof.Debug.Active {
			configRoutine = debugRoutineInterval
		}

		r := routines.NewRoutineJobs(prof, clientFactory, gh, done)

		mustSchedule(ctx, s.Every(updateRoutineInterval).StartImmediately().SingletonMode(), r.UpdateRoutine, prof.Name, "update")
		mustSchedule(ctx, s.Every(hzRoutineInterval).StartImmediately(), r.HZRoutine, prof.Name, "hz")
		mustSchedule(ctx, s.Every(configRoutine).StartImmediately().SingletonMode(), r.ConfigRoutine, prof.Name, "config")
		mustSchedule(ctx, s.Every(proxyRoutineInterval).StartImmediately(), r.ProxyRoutine, prof.Name, "proxy")
	}

	s.StartAsync()

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

// mustSchedule attaches the routine to the job being built on the scheduler and
// stops the agent when it cannot be scheduled, e.g. on an invalid interval,
// instead of dropping the job silently.
func mustSchedule(ctx context.Context, s *gocron.Scheduler, routine func(context.Context), profile, name string) {
	if _, err := s.Do(routine, ctx); err != nil {
		logrus.WithError(err).Fatalf("error while scheduling the %s routine for the profile %s", name, profile)
	}
}
