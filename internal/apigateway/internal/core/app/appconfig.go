package app

type Config struct {
	Server  ServerConfig
	FilmSrv ServiceConfig
	UserSrv ServiceConfig
}

type ServerConfig struct {
	Port uint16
}

type ServiceConfig struct {
	Host string
	Port uint16
}
