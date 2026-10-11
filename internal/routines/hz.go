package routines

import (
	"context"
	"fmt"
	"net"
	"os"

	"connectrpc.com/connect"
	agents_publicv1 "pkg.redcarbon.ai/proto/redcarbon/agents_public/v1"
)

func (r RoutineConfig) HZRoutine(ctx context.Context) {
	hostname, err := os.Hostname()
	if err != nil {
		r.Logger().Error("Error while retrieving the Hostname", "error", err)
		os.Exit(1)
	}

	localAddr, err := getOutboundIP()
	if err != nil {
		r.Logger().Error("Error while retrieving the IP address", "error", err)
		return
	}

	req := connect.NewRequest(&agents_publicv1.HZRequest{
		Ip:       localAddr,
		Hostname: hostname,
	})

	req.Header().Set("authorization", fmt.Sprintf("ApiToken %s", r.profile.Profile.Token))

	res, err := r.agentsCli.HZ(ctx, req)
	if err != nil {
		r.Logger().Error("Error while calling the HZ", "error", err)
		return
	}

	r.Logger().Info("HZ executed", "received_at", res.Msg.ReceivedAt.AsTime())
}

func getOutboundIP() (string, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "", err
	}

	defer conn.Close()

	return conn.LocalAddr().String(), nil
}
