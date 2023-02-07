package routines

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	"pkg.redcarbon.ai/internal/sentinelone"
	agentsExternalApiV1 "pkg.redcarbon.ai/proto/redcarbon/external_api/agents/api/v1"
)

func (r routineConfig) ConfigRoutine() {
	logrus.Infof("Start pulling the configurations from the server...\n")

	configs, err := r.agentsCli.PullConfigurations(context.Background(), &agentsExternalApiV1.PullConfigurationsReq{})
	if err != nil {
		logrus.Errorf("Error while retrieving the configurations for error %v", err)
		return
	}

	logrus.Infof("Configurations successfully pulled! Starting the jobs...")

	ctx, cFn := context.WithTimeout(context.Background(), time.Hour)

	var wg sync.WaitGroup

	for _, conf := range configs.AgentConfigurations {
		if conf.Data.GetSentinelOne() != nil {
			r.runService(ctx, sentinelone.RunSentinelOneService, &wg, conf)
			continue
		}
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

func (r routineConfig) runService(ctx context.Context, runner func(context.Context, *agentsExternalApiV1.AgentConfiguration, agentsExternalApiV1.AgentsExternalV1SrvClient), wg *sync.WaitGroup, ac *agentsExternalApiV1.AgentConfiguration) {
	wg.Add(1)

	go func() {
		runner(ctx, ac, r.agentsCli)
		wg.Done()
	}()
}
