package db

type DatabaseConfig struct {
	Host     string
	Port     uint16
	Name     string
	User     string
	Password string
	SslMode  string
}
