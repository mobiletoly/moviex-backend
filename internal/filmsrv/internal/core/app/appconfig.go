package app

import (
	"github.com/mobiletoly/moviex-backend/internal/common/dbc"
)

type Config struct {
	Server ServerConfig
	DB     dbc.DatabaseConfig
}

type ServerConfig struct {
	Port uint16
}
