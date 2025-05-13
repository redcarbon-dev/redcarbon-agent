package routines

import (
	"context"
	"fmt"
	"net"
	"os"

	"connectrpc.com/connect"
	"github.com/sirupsen/logrus"
	agents_publicv1 "pkg.redcarbon.ai/proto/redcarbon/agents_public/v1"
)

func (r RoutineConfig) HZRoutine(ctx context.Context) {
	hostname, err := os.Hostname()
	if err != nil {
		r.Logger().Fatalf("Error while retrieving the Hostname %v", err)
	}

	localAddr, err := getOutboundIP()
	if err != nil {
		return
	}

	req := connect.NewRequest(&agents_publicv1.HZRequest{
		Ip:       localAddr,
		Hostname: hostname,
	})

	req.Header().Set("authorization", fmt.Sprintf("ApiToken %s", r.profile.Profile.Token))

	res, err := r.agentsCli.HZ(ctx, req)
	if err != nil {
		logrus.WithError(err).Error("Error while calling the HZ")
		return
	}

	r.Logger().Infof("HZ executed %s", res.Msg.ReceivedAt.String())
}

func getOutboundIP() (string, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		logrus.WithError(err).Error("Error while retrieving the IP address")
		return "", err
	}

	defer conn.Close()

	return conn.LocalAddr().String(), nil
}
