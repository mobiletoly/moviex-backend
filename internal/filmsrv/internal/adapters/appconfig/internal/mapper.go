package internal

import (
	"github.com/mobiletoly/moviex-backend/internal/common/db"
)

func (cfg *RawAppConfig) ParseDBConfig() *db.DatabaseConfig {
	DB := cfg.Database
	return &db.DatabaseConfig{
		Host:     DB.Host,
		Port:     DB.Port,
		Name:     DB.Name,
		User:     DB.User,
		Password: DB.Password,
		SslMode:  DB.SslMode,
	}
}
