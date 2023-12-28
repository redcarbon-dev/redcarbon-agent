package routines

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"time"

	"connectrpc.com/connect"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"pkg.redcarbon.ai/internal/services"
	agents_publicv1 "pkg.redcarbon.ai/proto/redcarbon/agents_public/v1"
)

func (r RoutineConfig) ConfigRoutine(ctx context.Context) {
	logrus.Infof("Start pulling the configurations from the server...")

	ctxWithTimeout, cFn := context.WithTimeout(ctx, time.Hour)
	defer cFn()

	req := connect.NewRequest(&agents_publicv1.FetchAgentConfigurationRequest{})

	req.Header().Set("authorization", fmt.Sprintf("ApiToken %s", viper.Get("auth.access_token")))

	config, err := r.agentsCli.FetchAgentConfiguration(ctxWithTimeout, req)
	if err != nil {
		logrus.Errorf("Error while retrieving the configurations for error %v", err)
		return
	}

	logrus.Infof("Configurations successfully pulled! Starting the jobs...")

	var wg sync.WaitGroup

	jobs := services.NewServicesFromConfig(r.agentsCli, config.Msg.Configuration)

	for _, job := range jobs {
		r.runService(ctxWithTimeout, job, &wg)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		<-c
		cFn()
	}()

	wg.Wait()

	logrus.Infof("Jobs completed!")
}

func (r RoutineConfig) runService(ctx context.Context, s services.Service, wg *sync.WaitGroup) {
	wg.Add(1)

	go func() {
		s.RunService(ctx)
		wg.Done()
	}()
}
