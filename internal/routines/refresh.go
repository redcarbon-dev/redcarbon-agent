package routines

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func (r routineConfig) Refresh(ctx context.Context) {
	err := r.authSrv.RefreshToken(viper.GetString("auth.refresh_token"))
	if err != nil {
		logrus.Errorf("error while refreshing the access token %v", err)
	}
}
