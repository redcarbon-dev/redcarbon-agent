package auth

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc/metadata"

	agentsPublicApiV1 "pkg.redcarbon.ai/proto/redcarbon/public_apis/agents/api/v1"
)

type AuthenticationService struct {
	agentsCli  agentsPublicApiV1.AgentsPublicApiV1SrvClient
	configFile string
}

func NewAuthService(a agentsPublicApiV1.AgentsPublicApiV1SrvClient, configFile string) AuthenticationService {
	return AuthenticationService{
		agentsCli:  a,
		configFile: configFile,
	}
}

func (a AuthenticationService) RefreshToken(refreshToken string) error {
	ctx := context.Background()

	updateCtx := metadata.AppendToOutgoingContext(ctx, "authorization", fmt.Sprintf("Bearer %s", refreshToken))

	res, err := a.agentsCli.RefreshToken(updateCtx, &agentsPublicApiV1.RefreshTokenReq{})
	if err != nil {
		return err
	}

	viper.Set("auth.access_token", res.Token)
	viper.Set("auth.refresh_token", res.RefreshToken)

	err = viper.WriteConfigAs(a.configFile)
	if err != nil {
		logrus.Fatalf("Error while writing the configuration %v", err)
	}

	logrus.Infof("Token successfully refreshed!")

	return nil
}
