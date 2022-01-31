package appconfig

import (
	"github.com/mobiletoly/moviex-backend/internal/filmsrv/internal/adapters/appconfig/internal"
	"github.com/mobiletoly/moviex-backend/internal/filmsrv/internal/core/app"
)

func Load(deployment string) *app.Config {
	rawCfg := internal.LoadAppConfig(deployment)
	return &app.Config{
		Server: rawCfg.Server.ParseServerConfig(),
		DB:     rawCfg.Database.ParseDBConfig(),
	}
}
