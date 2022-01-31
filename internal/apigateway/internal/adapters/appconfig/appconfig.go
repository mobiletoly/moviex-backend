package appconfig

import (
	"github.com/mobiletoly/moviex-backend/internal/apigateway/internal/adapters/appconfig/internal"
	"github.com/mobiletoly/moviex-backend/internal/apigateway/internal/core/app"
)

func Load(deployment string) *app.Config {
	rawCfg := internal.LoadAppConfig(deployment)
	return &app.Config{
		Server:  rawCfg.Server.ParseServerConfig(),
		FilmSrv: rawCfg.Services.FilmSrv.ParseServiceConfig(),
		UserSrv: rawCfg.Services.UserSrv.ParseServiceConfig(),
	}
}
