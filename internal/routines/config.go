package routines

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc/metadata"

	"pkg.redcarbon.ai/internal/services"
	agentsExternalApiV1 "pkg.redcarbon.ai/proto/redcarbon/external_api/agents/api/v1"
)

func (r routineConfig) ConfigRoutine() {
	logrus.Infof("Start pulling the configurations from the server...\n")

	ctx := context.Background()

	ctxWithTimeout, cFn := context.WithTimeout(ctx, time.Hour)
	ctxWithTimeAndMeta := metadata.AppendToOutgoingContext(ctxWithTimeout, "authorization", fmt.Sprintf("Bearer %s", viper.Get("auth.access_token")))

	configs, err := r.agentsCli.PullConfigurations(ctxWithTimeAndMeta, &agentsExternalApiV1.PullConfigurationsReq{})
	if err != nil {
		logrus.Errorf("Error while retrieving the configurations for error %v", err)
		return
	}

	logrus.Infof("Configurations successfully pulled! Starting the jobs...")

	var wg sync.WaitGroup

	for _, conf := range configs.AgentConfigurations {
		s := services.NewServiceFromConfiguration(conf, r.agentsCli)
		if s == nil {
			logrus.Warnf("Configuration %s skipped as is not supported.", conf.AgentConfigurationId)
		}

		r.runService(ctxWithTimeAndMeta, s, &wg)
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

func (r routineConfig) runService(ctx context.Context, s services.Service, wg *sync.WaitGroup) {
	wg.Add(1)

	go func() {
		s.RunService(ctx)
		wg.Done()
	}()
}
