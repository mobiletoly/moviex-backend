package internal

import (
	"github.com/mobiletoly/moviex-backend/internal/common/dbc"
)

func (cfg *RawAppConfig) ParseDBConfig() *dbc.DatabaseConfig {
	DB := cfg.Database
	return &dbc.DatabaseConfig{
		Host:     DB.Host,
		Port:     DB.Port,
		Name:     DB.Name,
		User:     DB.User,
		Password: DB.Password,
		SslMode:  DB.SslMode,
	}
}
