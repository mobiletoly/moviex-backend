package internal

import (
	"github.com/mobiletoly/moviex-backend/internal/apigateway/internal/core/app"
)

func (cfg *RawAppServiceConfig) ParseServiceConfig() *app.ServiceConfig {
	return &app.ServiceConfig{
		Host: cfg.Host,
		Port: cfg.Port,
	}
}
