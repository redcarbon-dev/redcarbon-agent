package routines

import (
	"context"
	"github.com/sirupsen/logrus"
	"os"

	agentsExternalApiV1 "pkg.redcarbon.ai/proto/redcarbon/external_api/agents/api/v1"
)

func (r routineConfig) HZRoutine() {
	hostname, err := os.Hostname()
	if err != nil {
		logrus.Fatalf("Error while retrieving the Hostname %v", err)
	}

	res, err := r.agentsCli.HZ(context.Background(), &agentsExternalApiV1.HZReq{
		AgentId:    "agent:cld1k8f5g0000v6nm8xkc83b3",
		CustomerId: "c1",
		Ip:         "192.168.1.1",
		Hostname:   hostname,
	})
	if err != nil {
		// TODO Decide how to handle it
		logrus.Errorf("Error while calling the HZ %v", err)
		return
	}

	logrus.Printf("HZ executed %s", res.ReceivedAt.String())
}
