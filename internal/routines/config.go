package routines

import (
	"context"

	"github.com/sirupsen/logrus"

	agentsExternalApiV1 "pkg.redcarbon.ai/proto/redcarbon/external_api/agents/api/v1"
)

func (r routineConfig) ConfigRoutine() {
	_, err := r.agentsCli.PullConfigurations(context.Background(), &agentsExternalApiV1.PullConfigurationsReq{})
	if err != nil {
		logrus.Errorf("Error while retrieving the configurations for error %v", err)
		return
	}
}
