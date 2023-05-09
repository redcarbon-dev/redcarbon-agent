package routines

import (
	"context"
	"fmt"
	"net"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc/metadata"

	agentsPublicApiV1 "pkg.redcarbon.ai/proto/redcarbon/public_apis/agents/api/v1"
)

func (r routineConfig) HZRoutine(ctx context.Context) {
	hostname, err := os.Hostname()
	if err != nil {
		logrus.Fatalf("Error while retrieving the Hostname %v", err)
	}

	localAddr, err := getOutboundIP()
	if err != nil {
		return
	}

	updCtx := metadata.AppendToOutgoingContext(ctx, "authorization", fmt.Sprintf("Bearer %s", viper.Get("auth.access_token")))

	res, err := r.agentsCli.HZ(updCtx, &agentsPublicApiV1.HZReq{
		Ip:       localAddr,
		Hostname: hostname,
	})
	if err != nil {
		// TODO Decide how to handle it
		logrus.Errorf("Error while calling the HZ %v", err)
		return
	}

	logrus.Infof("HZ executed %s", res.ReceivedAt.String())
}

func getOutboundIP() (string, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		logrus.Errorf("Error while retrieving the IP address %v", err)
		return "", err
	}

	defer conn.Close()

	return conn.LocalAddr().String(), nil
}
