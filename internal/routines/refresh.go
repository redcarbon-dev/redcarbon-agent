package routines

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func (r routineConfig) Refresh() {
	err := r.authSrv.RefreshToken(viper.GetString("auth.refresh_token"))
	if err != nil {
		logrus.Errorf("error while refreshing the access token %v", err)
	}
}
