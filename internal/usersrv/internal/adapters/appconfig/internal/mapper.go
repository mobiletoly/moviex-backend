package internal

import (
	"github.com/mobiletoly/moviex-backend/internal/common/dbc"
	"github.com/mobiletoly/moviex-backend/internal/usersrv/internal/core/app"
)

func (cfg *RawAppDatabaseConfig) ParseDBConfig() dbc.DatabaseConfig {
	return dbc.DatabaseConfig{
		Host:     cfg.Host,
		Port:     cfg.Port,
		Name:     cfg.Name,
		User:     cfg.User,
		Password: cfg.Password,
		SslMode:  cfg.SslMode,
	}
}

func (cfg *RawAppServerConfig) ParseServerConfig() app.ServerConfig {
	return app.ServerConfig{
		Port: cfg.Port,
	}
}
