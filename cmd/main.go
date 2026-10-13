package main

import (
	"context"
	"log/slog"
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

	"github.com/spf13/cobra"

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
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, nil)))

	confDir, err := os.Getwd()
	if err != nil {
		slog.Error("can't extract the user working directory", "error", err)
		os.Exit(1)
	}

	if err = config.Init(confDir); err != nil {
		slog.Error("error while initializing the configuration", "error", err)
		os.Exit(1)
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
	slog.Info("Starting RedCarbon Agent", "version", build.Version, "arch", build.Architecture)

	conf := config.LoadConfiguration()

	if len(conf.Profiles) == 0 {
		slog.Error("No profiles found, please add one by running `redcarbon profile add`")
		os.Exit(1)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT)
	defer cancel()

	g, ctx := errgroup.WithContext(ctx)

	clientFactory := cli.NewClientFactory()
	gh := github.NewClient(nil, build.GitHubAPIURL)
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
		slog.Error("error while scheduling the routines", "error", err)
		os.Exit(1)
	}

	g.Go(func() error {
		<-ctx.Done()
		s.Stop()
		slog.Info("RedCarbon Agent stopped")
		os.Exit(0)
		return nil
	})

	g.Go(func() error {
		<-done
		s.Stop()
		slog.Info("RedCarbon Agent stopped due update")
		os.Exit(updateErrorCode)
		return nil
	})

	if err := g.Wait(); err != nil {
		slog.Error("error while running the agent", "error", err)
		os.Exit(1)
	}
}
