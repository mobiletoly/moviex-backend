package appconfig

import (
	"github.com/mobiletoly/moviex-backend/internal/apigateway/internal/adapters/appconfig/internal"
	"github.com/mobiletoly/moviex-backend/internal/apigateway/internal/core/app"
)

func Load() *app.Config {
	rawCfg := internal.LoadAppConfig()
	return &app.Config{
		FilmSrv: rawCfg.Services.FilmSrv.ParseServiceConfig(),
		UserSrv: rawCfg.Services.UserSrv.ParseServiceConfig(),
	}
}
