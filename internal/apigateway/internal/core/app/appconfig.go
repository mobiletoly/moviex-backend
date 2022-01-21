package app

type Config struct {
	FilmSrv *ServiceConfig
	UserSrv *ServiceConfig
}

type ServiceConfig struct {
	Host string
	Port uint16
}
