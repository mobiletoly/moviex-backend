package appconfig

import (
	"github.com/mobiletoly/moviex-backend/internal/usersrv/internal/adapters/appconfig/internal"
	"github.com/mobiletoly/moviex-backend/internal/usersrv/internal/core/app"
)

func Load() *app.Config {
	rawCfg := internal.LoadAppConfig()
	return &app.Config{
		DB: rawCfg.ParseDBConfig(),
	}
}
